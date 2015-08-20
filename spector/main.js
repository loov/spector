depends("/spector/trace.js");
depends("/spector/trace.view.js");
depends("/spector/import/simulator.js");

var stream = new spector.import.simulator.Stream();
var trace = new spector.Trace();
var traceview = new spector.Trace.View();

var canvas = document.getElementById("view");
var context = canvas.getContext("2d");

function render(){
	canvas.width = window.innerWidth - 4;
	canvas.height = window.innerHeight - 4;
	var size = {x: canvas.width, y: canvas.height};
	context.clearRect(0, 0, size.x, size.y);

	trace.consume(stream);
	traceview.render(context, trace, size);
	requestAnimationFrame(render);
}

render();