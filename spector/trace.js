package("spector", function(){
	depends("/spector/stream.js");
	depends("/spector/protocol.js");

	"use strict";
	// Trace is the ViewModel for main display
	var ValueArray = Array;
	var TimeArray  = Array;
	var IDArray    = Array;


	function Layer(){
		this.time_ = new TimeArray();
		this.id_ = new IDArray();
	}

	Layer.prototype = {
		get timepairs(){ return this.time_; },
		get IDs(){ return this.id_; },
		get lastTime(){
			if(this.time_.length == 0){ return -Infinity; }
			return this.time_[this.time_.length-1];
		},
		get lastID(){
			return this.id_[this.id_.length-1];
		},

		beginEvent: function(id, time){
			assert(this.lastTime <= time);
			this.time_.push(time);
			this.time_.push(Infinity);
			this.id_.push(id);
		},
		endEvent: function(id, time){
			assert(this.lastID === id);
			this.time_[this.time_.length-1] = time;
		}
	};

	function Thread(threadID, time){
		this.ThreadID  = threadID;
		this.Time = time;

		this.start = time;
		this.stop  = undefined;
		this.layers = [];
	}
	Thread.prototype = {
		getOpenLayer: function(time){
			var layers = this.layers;
			for(var i = 0; i < layers.length; i++){
				if(layers[i].lastTime <= time){
					return layers[i];
				}
			}

			var layer = new Layer();
			layers.push(layer);
			return layer;
		},
		getLayerWith: function(id){
			var layers = this.layers;
			for(var i = 0; i < layers.length; i++){
				if(layers[i].lastID === id){
					return layers[i];
				}
			}
			return undefined;
		},
		stopLayers: function(time){
			var layers = this.layers;
			for(var i = 0; i < layers.length; i++){
				if(!isFinite(layers[i].lastTime)){
					layers[i].endEvent(layers[i].lastID, time);
				}
			}
		}
	};

	function Track(){
		this.threads = [];
	}
	Track.prototype = {
		get first(){ return this.threads[0]; },
		get last(){ return this.threads[this.threads.length-1]; },
		get start(){
			var first = this.first;
			if(first === undefined) { return Infinity; }
			return first.start;
		},
		get stop(){
			var last = this.last;
			if(last === undefined){ return -Infinity; }
			if(last.stop === undefined){ return Infinity; }
			return last.stop;
		}
	};

	function Process(processID, machineID, time, cpufreq){
		this.ProcessID = processID;
		this.MachineID = machineID;
		this.Time = time;
		this.CPUFrequency = cpufreq;

		this.start = time;
		this.stop  = undefined;

		this.threads = new Array();
		this.tracks  = new Array();
	}
	Process.prototype = {
		getOpenTrack: function(time){
			var tracks = this.tracks;
			for(var i = 0; i < tracks.length; i++){
				if(tracks[i].stop <= time) {
					return tracks[i];
				}
			}

			var track = new Track();
			tracks.push(track);
			return track;
		},
		threadByID: function(threadID) {
			var threads = this.threads;
			for(var i = 0; i < threads.length; i++){
				if(threads[i].ThreadID == threadID){
					return threads[i];
				}
			}
			return undefined;
		}
	};

	function nowms() { return (new Date()).getTime(); }

	function Trace(){
		this.processes = new Array();
		this.totalEventCount = 0;
	}

	Trace.prototype = {
		get begin(){
			var procs = this.processes;
			if(procs.length === 0){ return 0; }

			var start = procs[0].start;
			for(var i = 1; i < procs.length; i+=1){
				start = Math.min(start, procs[i].start);
			}
			return start;
		},
		get end(){
			var procs = this.processes;
			if(procs.length === 0){ return 0; }

			var end = procs[0].Time;
			for(var i = 1; i < procs.length; i+=1){
				end = Math.max(end, procs[i].Time);
			}
			return end;
		},

		processByID: function(pid, mid) {
			var procs = this.processes;
			for(var i = 0; i < procs.length; i++){
				if((procs[i].ProcessID == pid) && (procs[i].MachineID == mid)){
					return procs[i];
				}
			}
			return undefined;
		},

		consume: function(stream, maxMillis) {
			maxMillis = maxMillis || 1000;
			var start = nowms();
			main:
			while(nowms() - start < maxMillis){
				for(var i = 0; i < 8; i++){
					var ev = stream.next();
					if(ev === undefined){
						break main;
					}

					var fn = accept[ev.Code];
					if(fn === undefined){
						console.log("unhandled", ev.Code, ev);
						continue;
					}

					this.totalEventCount += 1;

					var cache = stream.meta;
					if(cache.TraceProcess !== undefined){
						if(ev.Time !== undefined){
							cache.TraceProcess.Time = ev.Time;
						}
					}

					fn(cache, this, ev);
				}
			}
		}
	};

	function assert(ok, msg){
		if(!ok){ throw new Error(msg); }
	}

	var accept = {};
	function when(code, fn){ accept[code] = fn; }

	var code = spector.EventCode;
	when(code.StreamStart, function(cache, trace, event){
		var proc = new Process(
			event.ProcessID,
			event.MachineID,
			event.Time,
			event.CPUFrequency
		);
		cache.TraceProcess = proc;

		trace.processes.push(proc);
	});

	when(code.StreamStop, function(cache, trace, event){
		cache.TraceProcess = undefined;
	});

	when(code.ThreadStart, function(cache, trace, event){
		var proc = cache.TraceProcess;
		assert(proc !== undefined);

		var thread = new Thread(event.ThreadID, event.Time);
		var track = proc.getOpenTrack(event.Time);
		track.threads.push(thread);
		proc.threads.push(thread);
	});

	when(code.ThreadStop, function(cache, trace, event){
		var proc = cache.TraceProcess;
		assert(proc !== undefined);

		var thread = proc.threadByID(event.ThreadID);
		thread.Time = event.Time;
		thread.stop = event.Time;
		thread.stopLayers(event.Time);
	});

	when(code.ThreadSleep, function(cache, trace, event){
		var proc = cache.TraceProcess;
		assert(proc !== undefined);
	});

	when(code.ThreadWake, function(cache, trace, event){
		var proc = cache.TraceProcess;
		assert(proc !== undefined);
	});


	when(code.Begin, function(cache, trace, event){
		var proc = cache.TraceProcess;
		assert(proc !== undefined);

		var thread = proc.threadByID(event.ThreadID);
		thread.Time = event.Time;

		var layer = thread.getOpenLayer(event.Time);
		layer.beginEvent(event.ID, event.Time);
	});

	when(code.End, function(cache, trace, event){
		var proc = cache.TraceProcess;
		assert(proc !== undefined);

		var thread = proc.threadByID(event.ThreadID);
		thread.Time = event.Time;

		var layer = thread.getLayerWith(event.ID);
		layer.endEvent(event.ID, event.Time);
	});

	return {
		Trace: Trace
	};
});