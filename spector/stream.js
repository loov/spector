package("spector", function(){
	depends("/spector/util/utf8.js");
	depends("/spector/protocol.js");

	var Magic = "spector";

	var Stage = {
		Header:  0, // waiting for header
		Reading: 1, // reading events
		Stopped: 2, // after StreamStop
	};

	Stream.Stage = Stage;
	function Stream(bufferSize) {
		// should we use here a list of Uint8Array blocks instead?
		this.buffer = new Uint8Array(bufferSize);
		this.len  = 0;
		this.head = 0; // reading head
		this.version = spector.Version;
		this.stage = Stage.New;

		this.eventPool_ = new Array();
	}

	Stream.prototype = {
		collect_: function(){
			if (this.head <= 0) { return; }
			this.buffer.set(0, this.buffer.slice(this.head, this.len));
			this.len = this.len - this.head;
			this.head = 0;
		},
		realloc_: function(newSize){
			var nextSize = (newSize*3/2)|0;
			var buffer = new Uint8Array(newSize);
			buffer.set(0, this.buffer.slice(this.head, this.len));
			this.len = this.len - this.head;
			this.head = 0;
			this.buffer = buffer;
		},

		// partial content must not be appended!
		write_: function(data){
			if(this.stage >= Stage.Stopped){
				throw new Error("stream has been stopped and no more data can be added!");
			}
			if (this.len + data.length > this.buffer.length) {
				if(this.len - this.head + data.length > this.buffer.length) {
					this.collect_();
				} else {
					this.realloc_(this.len + data.length);
				}
			}
			this.buffer.set(this.len, data);
		},

		read_: function(len){
			if(this.head + len > this.len){
				throw new Error("not enough data!");
			}
			var content = this.buffer.slice(this.head, this.head + len);
			this.head += len;
			return content;
		},
		readMagic_: function(){
			var magic = this.read_(Magic.length);
			for(var i = 0; i < magic.length; i++){
				if(magic[i] !== Magic.charCodeAt(i)){
					throw new Error("invalid magic header: " +
						String.fromCharCode.apply("", magic));
				}
			}

			this.version = this.readInt();
			this.stage = Stage.Reading;
		},

		readByte: function(){
			var val = this.buffer[this.head];
			this.head += 1;
			return val;
		},
		// TODO: use some other encoding for integers, it's a poor choice
		readInt: function(){
			var u8 = this.buffer.slice(this.head, this.head+4);
			var uint = new Int32Array(u8);
			var val = uint[0];
			this.head += 4;
			return val;
		},
		readBlob: function(){
			var count = this.readInt();
			return this.read_(count);
		},
		readUTF8: function(){
			var count = this.readInt();
			var bytes = this.read_(count);
			return spector.util.utf8.decode(bytes);
		},
		readValues: function(){
			var count = this.readInt();
			var result = new Array(count*2);
			for(var i = 0; i < count; i++){
				result[i*2] = this.readInt();
				result[i*2+1] = this.readInt();
			}
			return result;
		},

		// !!! the returned object will be reused                 !!!
		// !!! the values must be copied to the final destination !!!
		next: function(){
			if (this.stage === Stage.Header) {
				if(this.len - this.head <= Magic.length)
					return false;
				this.readMagic_();
			} else if (this.stage === Stage.Stopped) {
				throw new Error("stream has been stopped!");
			}

			if(this.head >= this.len){
				return undefined;
			}

			var EventCode = this.readByte();
			var obj = this.eventPool_[EventCode];
			if(obj === undefined){
				obj = new spector.EventByCode[EventCode](this);
				this.eventPool_[EventCode] = obj;
			}
			obj.read(this);
			if(obj.Type === spector.Event.StreamStop.Type){
				this.stage = Stage.Stopped;
			}

			return obj;
		}
	};

	function Value(id, value) {
		this.ID = id;
		this.Value = value;
	}

	return {
		Stream: Stream
	};
});