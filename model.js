function assert(v, msg){
	if(!v){
		throw new Error(msg);
	}
}

function find(values, val){
	return values.indexOf(val);
};

var ValueArray = Array;
var TimeArray  = Array; // Float64Array;
var EIDArray   = Array; // Uint32Array;
var IDArray    = Array; // Uint32Array;

function Durations(){
	this.time_ = new TimeArray();
	this.eids_ = new EIDArray();
	this.openName_ = undefined;
	this.view = {};
}

Durations.prototype = {
	get count(){ return this.time_.length >> 1; },
	begin: function(i){ return this.time_[i<<1]; },
	end:   function(i){ return this.time_[(i<<1) + 1]; },
	eid:    function(i){ return this.eids_[i]; },

	addBegin: function(eid, name, ts){
		assert(this.openName_ === undefined, "already open");
		assert(this.mark < ts, "invalid order");
		this.time_.push(ts);
		this.eids_.push(eid);
		this.openName_ = name;
	},

	addEnd: function(eid, name, ts){
		assert(this.mark < ts, "invalid order");
		assert(this.openName_ === name, "already open");
		this.time_.push(ts);
		this.eids_.push(eid);
		this.openName_ = undefined;
	},

	get openName(){
		return this.openName_;
	},
	get closedAt(){
		return this.openName_ === undefined ? this.mark : Infinity;
	},
	get mark(){
		if(this.time_.length === 0){ return -Infinity; }
		return this.time_[this.time_.length-1];
	}
};

function DurationsStack(){
	this.layers = [];
	this.view = {};
}

DurationsStack.prototype = {
	addBegin: function(eid, name, ts){
		for(var i = 0; i < this.layers.length; i++){
			var layer = this.layers[i];
			if(layer.closedAt <= ts){
				layer.addBegin(eid, name, ts);
				return;
			}
		}
		var dur = new Durations();
		dur.addBegin(eid, name, ts);
		this.layers.push(dur);
	},
	addEnd: function(eid, name, ts){
		for(var i = 0; i < this.layers.length; i++){
			var layer = this.layers[i];
			if(layer.openName == name){
				this.layers[i].addEnd(eid, name, ts);
				return;
			}
		}
		assert(false, "beginning not found");
	}
};

function Counter(name){
	this.name = name;
	this.values_ = new ValueArray();
	this.time_   = new TimeArray();
	this.eids_   = new EIDArray();

	this.view = {};
}

Counter.prototype = {
	add: function(eid, ts, value){
		this.eids_.push(eid);
		this.time_.push(ts);
		this.values_.push(value);
	},

	get count(){ return this.time_.length; },
	eid: function(i){ return this.eids_[i]; },
	time: function(i){ return this.time_[i]; },
	value: function(i){ return this.value_[i]; }
};

function CounterStack(){
	this.counters = [];
	this.byName_ = {};

	this.view = {};
}

CounterStack.prototype = {
	add: function(eid, name, ts, value){
		var counter = this.byName_[name];
		if(counter === undefined){
			counter = new Counter(name);
			this.counters.push(counter);
			this.byName_[name] = counter;
		}

		counter.add(eid, ts, value);
	},
};

function InstantsStack(){
	this.ids_ = new EIDArray();
	this.time_ = new TimeArray();
};

InstantsStack.prototype = {
	get count(){ return this.ids_.length; },
	time: function(i){ return this.time_[i]; },
	add: function(eid, ts){
		this.ids_.push(eid);
		this.time_.push(ts);
	}
}

function Thread(pid, tid){
	this.pid = pid;
	this.tid = tid;

	this.durations = new DurationsStack();
	this.counters  = new CounterStack();
	this.instants = new InstantsStack();

	this.view = {};
}

Thread.prototype = {
	add: function(eid, ev){
		var ph = ev.ph;
		if (ph == "B") {
			this.durations.addBegin(eid, ev.name, ev.ts);
		} else if (ph === "E") {
			this.durations.addEnd(eid, ev.name, ev.ts);
		} else if (ph === "X") {
			this.durations.addBegin(eid, ev.name, ev.ts);
			this.durations.addEnd(eid, ev.name, ev.ts + ev.dur);
		} else if (ph === "C") {
			for(var name in ev.args){
				if(ev.args.hasOwnProperty(name)){
					var arg = ev.args[name];
					this.counters.add(eid, ev.name, ev.ts, arg);
				}
			}
		} else if (ph === "I") {
			this.instants.add(eid, ev.ts);
		} else {
			console.log("unimplemented", ph);
		}
	}
};

function Process(pid){
	this.pid = pid;
	this.threads = [];
	this.instants = new InstantsStack();

	this.view = {};
}

Process.prototype = {
	threadByID: function(tid){
		for(var i = 0; i < this.threads.length; i++){
			var thread = this.threads[i];
			if(thread.tid == tid){
				return thread;
			}
		}
		var thread = new Thread(this.pid, tid);
		this.threads.push(thread);
		this.threads.sort(function(a,b){ return a.tid - b.tid; });
		return thread;
	},
	add: function(eid, ev){
		assert(this.pid == ev.pid);
		if((ev.ph === "I") && (ev.s == "p")){
			this.instants.push(eid, ev.ts);
			return;
		}

		this.threadByID(ev.tid).add(eid, ev);
	}
};

function FlowEnd(eid, ev, layer, ts){
	this.id  = ev.id;

	this.pid = ev.pid;
	this.tid = ev.tid;
	this.layer = layer;

	this.ts    = ts;
}

function Flows(){
	this.ids_   = new IDArray();

	this.startEID_  = new EIDArray();
	this.start_     = new TimeArray();

	this.finishEID_ = new EIDArray();
	this.finish_    = new TimeArray();
}

Flows.prototype = {
	get count(){ return this.ids_.length; },
	id: function(i){ return this.ids_[i]; },
	start: function(i){ return this.start_[i]; },
	startEID: function(i){ return this.startEID_[i]; },
	start: function(i){ return this.start_[i]; },
	finishEID: function(i){ return this.finishEID_[i]; },
	finish: function(i){ return this.finish_[i]; },

	addStart: function(eid, id, ts){
		this.ids_.push(id);
		this.startEID_.push(eid);
		this.start_.push(ts);
		this.finishEID_.push(0);
		this.finish_.push(0);
	},
	addFinish: function(eid, id, ts){
		var i = find(this.ids_, id);
		assert(i >= 0, "missing id " + id);
		this.finishEID_[i] = eid;
		this.finish_[i] = ts;
	}
};

function Timeline(){
	this.events = [];
	this.processes = [];
	this.instants = new InstantsStack();
	this.flows = new Flows();

	this.view = {};
}

Timeline.prototype = {
	processByID: function(pid){
		for(var i = 0; i < this.processes.length; i++){
			var process = this.processes[i];
			if(process.pid == pid){
				return process;
			}
		}
		var process = new Process(pid);
		this.processes.push(process);
		this.processes.sort(function(a,b){ return a.pid - b.pid; });
		return process;
	},
	trackByID: function(pid, tid){
		return this.processByID(pid).threadByID(tid);
	},
	add: function(ev){
		try {
			if(ev.ph === "M"){ return; }
			var eid = this.events.length;
			this.events.push(ev);

			if((ev.ph === "I") && (ev.s == "g")){
				this.instants.push(eid, ev.ts);
			} else if (ev.ph === "s") {
				//TODO: find enclosing slice end
				this.flows.addStart(eid, ev.id, ev.ts);
			} else if (ev.ph === "t") {
				//TODO: find enclosing slice start
				this.flows.addFinish(eid, ev.id, ev.ts);
			} else if (ev.ph === "f"){
				//TODO: find next slice
				// if ev.bp == e, find enclosing slice start
				this.flows.addFinish(eid, ev.id, ev.ts);
			} else {
				this.processByID(ev.pid).add(eid, ev);
			}
		} catch(e){
			console.error(e);
		}
	},

	load: function(model){
		for(var i = 0; i < model.traceEvents.length; i++){
			var ev = model.traceEvents[i];
			this.add(ev);
		}

		this.zoomToPage();
	},

	zoomToPage: function(){
		var v = this.view;
		v.start = Infinity;
		v.end = -Infinity;
		this.events.map(function(ev){
			v.start = Math.min(v.start, ev.ts);
			v.end = Math.max(v.end, ev.ts + (ev.dur || 0));
		});
	},

	// 1 to zoom in, -1 to zoom out by a factor of 2
	zoom: function(px, amount){
		var v = this.view;
		var center = v.start + (v.end - v.start) * px;
		var scale = Math.pow(2, -amount);

		v.start = center + (v.start - center) * scale;
		v.end = center + (v.end - center) * scale;
	},

	// dx is normalized offset in x
	drag: function(dx){
		var v = this.view;
		var dt = (v.end - v.start) * dx;
		v.start += dt;
		v.end += dt;
	}
};