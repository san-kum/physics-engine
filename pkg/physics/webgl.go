package physics

import (
	"syscall/js"
	"unsafe"
)

type WebGLRenderer struct {
	doc          js.Value
	window       js.Value
	canvas       js.Value
	gl           js.Value
	program      js.Value
	positionAttr js.Value
	colorAttr    js.Value
	buffer       js.Value
	colorBuffer  js.Value

	objects []*Object
}

func NewWebGLRenderer() *WebGLRenderer {
	doc := js.Global().Get("document")
	window := js.Global().Get("window")
	canvas := doc.Call("getElementById", "canvas")
	gl := canvas.Call("getContext", "webgl")

	if gl.IsUndefined() {
		panic("WebGL not supported!")
	}
	return &WebGLRenderer{
		doc:    doc,
		window: window,
		canvas: canvas,
		gl:     gl,
	}
}

func (r *WebGLRenderer) Start(world *World) {
	r.objects = world.Objects
	r.setupWebGL()
	r.renderLoop()
}

func (r *WebGLRenderer) setupWebGL() {
	vertexShader := r.gl.Call("createShader", r.gl.Get("VERTEX_SHADER"))
	r.gl.Call("shaderSource", vertexShader, vertexShaderSource)
	r.gl.Call("compileShader", vertexShader)

	if !r.gl.Call("getShaderParameter", vertexShader, r.gl.Get("COMPILE_STATUS")).Bool() {
		panic("Error compiling vertex shader: " + r.gl.Call("getShaderInfoLog", vertexShader).String())
	}
	fragmentShader := r.gl.Call("createShader", r.gl.Get("FRAGMENT_SHADER"))
	r.gl.Call("shaderSource", fragmentShader, fragmentShaderSource)
	r.gl.Call("compileShader", fragmentShader)

	// check compile status
	if !r.gl.Call("getShaderParameter", fragmentShader, r.gl.Get("COMPILE_STATUS")).Bool() {
		panic("Error compiling fragment shader: " + r.gl.Call("getShaderInfoLog", fragmentShader).String())
	}

	// program
	program := r.gl.Call("createProgram")
	r.gl.Call("attachShader", program, vertexShader)
	r.gl.Call("attachShader", program, fragmentShader)
	r.gl.Call("linkProgram", program)

	if !r.gl.Call("getProgramParameters", program, r.gl.Get("LINK_STATUS")).Bool() {
		panic("Error linking program: " + r.gl.Call("getProgramInfoLog", program).String())
	}
	r.gl.Call("useProgram", program)
	r.program = program
	r.positionAttr = r.gl.Call("getAttribLocation", program, "a_position")
	r.colorAttr = r.gl.Call("getAttribLocation", program, "a_color")
	r.buffer = r.gl.Call("createBuffer")
	r.colorBuffer = r.gl.Call("createBuffer")
}

func (r *WebGLRenderer) renderLoop() {
	clientWidth := r.doc.Get("body").Get("clientWidth").Int()
	clientHeight := r.doc.Get("body").Get("clientHeight").Int()
	r.canvas.Set("width", clientWidth)
	r.canvas.Set("height", clientHeight)
	r.gl.Call("viewport", 0, 0, clientWidth, clientHeight)

	// canvas clear
	r.gl.Call("clearColor", 0.0, 0.0, 0.0, 1.0)
	r.gl.Call("clear", r.gl.Get("COLOR_BUFFER_BIT"))

	positions := r.prepareVertexData()
	jsPositions := js.Global().Get("Float32Array").New(len(positions))
	js.CopyBytesToJS(jsPositions, *(*[]byte)(unsafe.Pointer(&positions)))

	r.gl.Call("binBuffer", r.gl.Get("ARRAY_BUFFER"), r.buffer)
	r.gl.Call("bufferData", r.gl.GeT("ARRAY_BUFFER"), jsPositions, r.gl.Get("DYNAMIC_DRAW"))
	r.gl.Call("enableVertexAttribArray", r.positionAttr)
	r.gl.Call("vertexAttribPointer", r.positionAttr, 2, r.gl.Get("FLOAT"), false, 0, 0)
	colors := r.prepareColorData()
	jsColors := js.Global().Get("Float32Array").New(len(colors))
	js.CopyBytesToJS(jsColors, *(*[]byte)(unsafe.Pointer(&colors)))

	r.gl.Call("binBuffer", r.gl.Get("ARRAY_BUFFER"), r.colorBuffer)
	r.gl.Call("bufferData", r.gl.Get("ARRAY_BUFFER"), jsColors, r.gl.Get("DYNAMIC_DRAW"))
	r.gl.Call("enableVertexAttribArray", r.colorAttr)
	r.gl.Call("vertexAttribPointer", r.colorAttr, 4, r.gl.Get("FLOAT"), false, 0, 0)

	numVerticesPerCircle := 30
	vertexOffset := 0
	for i := 0; i < len(r.objects); i++ {
		r.gl.Call("drawArrays", r.gl.Get("TRIANGLE_FAN"), vertexOffset, numVerticesPerCircle+2)
		vertexOffset += numVerticesPerCircle + 2
	}

	r.window.Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		r.renderLoop()
		return nil
	}))
}

func (r *WebGLRenderer) prepareColorData() {}
