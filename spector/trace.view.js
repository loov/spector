package("spector.Trace", function(){
	"use strict";

	function UI(context, size, trace){
		this.context = context;
		this.size = size;
		this.Y = 0;

		this.begin = trace.begin;
		this.end   = trace.end;
		this.timespan  = this.end - this.begin;
		this.px    = this.timespan / size.x;
	}

	UI.prototype = {
		font: function(fontName) {
			this.context.font = fontName;
		},
		h1: function(text, height, color, background){
			var ctx = this.context;

			ctx.fillStyle = background;
			ctx.fillRect(0, this.Y, this.size.x, height);
			this.Y += height;

			ctx.fillStyle = color;
			ctx.font = height + "px";
			ctx.fillText(text, 4, this.Y - 5);
		},
		hr: function(color, height){
			var ctx = this.context;

			height = height || 1;
			ctx.fillStyle = color;
			ctx.fillRect(0, this.Y, this.size.x, height);
			this.Y += height;
		},
		pad: function(height){
			this.Y += height;
		},

		trX: function(x){
			if(x == undefined){ return this.size.x; }
			var v = ((x - this.begin) * this.size.x / this.timespan);
			if(v > this.size.x){ v = this.size.x; }
			return v | 0;
		},

		block: function(id, start, stop, height){
			var ctx = this.context;

			var x0 = this.trX(start);
			var x1 = this.trX(stop);

			this.Y += 10;
			ctx.fillStyle = "hsla(" + ((id * 34)%360) + ",50%,50%,1)";
			ctx.font = "8px";
			ctx.fillText(id, x0 + 2, this.Y);
			this.Y += 2;

			ctx.fillStyle = "hsla(" + ((id * 34)%360) + ",50%,90%,1)";
			ctx.fillRect(x0, this.Y, x1 - x0, height);
		},

		spans: function(ids, pairs, depth, height){
			var ctx = this.context;
			ctx.fillStyle = "hsla(0,40%," + (60 + (depth*20)|0) + "%,1)";

			var mingap = this.px;
			for(var i = 0; i < pairs.length; i += 2){
				if(pairs[i+1] < this.begin){ continue; }
				if(pairs[i] > this.end){ break; }

				var start = pairs[i];
				// connect small gaps
				while(pairs[i+2] - pairs[i+1] < mingap){
					i += 2;
				}
				var stop = pairs[i+1];

				var x0 = this.trX(start);
				var x1 = this.trX(stop);
				ctx.fillRect(x0, this.Y, x1 - x0, height);
			}
			this.Y += height;
		}
	};


	function View(){
		this.lastEnd = 0;
	}

	View.prototype = {
		render: function(ctx, trace, size){
			var ui = new UI(ctx, size, trace);
			this.lastEnd = (this.lastEnd + ui.end) / 2;
			// modify end
			ui.end = this.lastEnd;
			ui.timespan  = ui.end - ui.begin;
			ui.px    = ui.timespan / size.x;

			ui.font("Courier New");
			ui.h1(trace.totalEventCount, 18, "#fff", "#000");

			for(var pi = 0; pi < trace.processes.length; pi += 1){
				var proc = trace.processes[pi];
				ui.h1(proc.MachineID + " > " + proc.ProcessID, 18, "#333", "#eee");

				for(var ti = 0; ti < proc.tracks.length; ti += 1){
					var track = proc.tracks[ti];
					ui.pad(4);

					var MaximumY = ui.Y;
					var BaseY = ui.Y;

					const LayerHeight = 10;

					for(var thi = 0; thi < track.threads.length; thi += 1){
						ui.Y = BaseY;
						var thread = track.threads[thi];

						ui.block(
							thread.ThreadID,
							thread.start, thread.stop,
							thread.layers.length*LayerHeight
						);

						var layerlen = thread.layers.length;
						for(var li = 0; li < layerlen; li += 1){
							var layer = thread.layers[li];
							ui.spans(layer.IDs, layer.timepairs, (layerlen-li-1)/layerlen, LayerHeight);
						}
						ui.pad(3);

						MaximumY = Math.max(MaximumY, ui.Y);
					}
					ui.Y = MaximumY;
				}
			}
		}
	};

	return {
		View: View
	};
});