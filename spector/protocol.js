package("spector", function(){
	"use strict";

	var Event = {};
	var EventByCode = {};
	var EventCode = {};


	Event.Invalid = InvalidEvent;
	InvalidEvent.Code = 0x00;
	EventCode.Invalid = 0x00;
	EventByCode[0x00] = InvalidEvent;
	function InvalidEvent(props){
		props = props !== undefined ? props : {};

	};

	InvalidEvent.prototype = {
		Code: 0x00,
		read: function(stream){
		},
		write: function(stream){
		}
	};

	Event.StreamStart = StreamStartEvent;
	StreamStartEvent.Code = 0x01;
	EventCode.StreamStart = 0x01;
	EventByCode[0x01] = StreamStartEvent;
	function StreamStartEvent(props){
		props = props !== undefined ? props : {};

		this.ProcessID = props.ProcessID || 0;
		this.MachineID = props.MachineID || 0;
		this.Time = props.Time || 0;
		this.CPUFrequency = props.CPUFrequency || 0;
	};

	StreamStartEvent.prototype = {
		Code: 0x01,
		read: function(stream){
			this.ProcessID = stream.readInt();
			this.MachineID = stream.readInt();
			this.Time = stream.readInt();
			this.CPUFrequency = stream.readInt();
		},
		write: function(stream){
			stream.writeInt(this.ProcessID);
			stream.writeInt(this.MachineID);
			stream.writeInt(this.Time);
			stream.writeInt(this.CPUFrequency);
		}
	};

	Event.StreamStop = StreamStopEvent;
	StreamStopEvent.Code = 0x02;
	EventCode.StreamStop = 0x02;
	EventByCode[0x02] = StreamStopEvent;
	function StreamStopEvent(props){
		props = props !== undefined ? props : {};

		this.Time = props.Time || 0;
	};

	StreamStopEvent.prototype = {
		Code: 0x02,
		read: function(stream){
			this.Time = stream.readInt();
		},
		write: function(stream){
			stream.writeInt(this.Time);
		}
	};

	Event.ThreadStart = ThreadStartEvent;
	ThreadStartEvent.Code = 0x03;
	EventCode.ThreadStart = 0x03;
	EventByCode[0x03] = ThreadStartEvent;
	function ThreadStartEvent(props){
		props = props !== undefined ? props : {};

		this.Time = props.Time || 0;
		this.ThreadID = props.ThreadID || 0;
		this.StackID = props.StackID || 0;
	};

	ThreadStartEvent.prototype = {
		Code: 0x03,
		read: function(stream){
			this.Time = stream.readInt();
			this.ThreadID = stream.readInt();
			this.StackID = stream.readInt();
		},
		write: function(stream){
			stream.writeInt(this.Time);
			stream.writeInt(this.ThreadID);
			stream.writeInt(this.StackID);
		}
	};

	Event.ThreadSleep = ThreadSleepEvent;
	ThreadSleepEvent.Code = 0x04;
	EventCode.ThreadSleep = 0x04;
	EventByCode[0x04] = ThreadSleepEvent;
	function ThreadSleepEvent(props){
		props = props !== undefined ? props : {};

		this.Time = props.Time || 0;
		this.ThreadID = props.ThreadID || 0;
		this.StackID = props.StackID || 0;
	};

	ThreadSleepEvent.prototype = {
		Code: 0x04,
		read: function(stream){
			this.Time = stream.readInt();
			this.ThreadID = stream.readInt();
			this.StackID = stream.readInt();
		},
		write: function(stream){
			stream.writeInt(this.Time);
			stream.writeInt(this.ThreadID);
			stream.writeInt(this.StackID);
		}
	};

	Event.ThreadWake = ThreadWakeEvent;
	ThreadWakeEvent.Code = 0x05;
	EventCode.ThreadWake = 0x05;
	EventByCode[0x05] = ThreadWakeEvent;
	function ThreadWakeEvent(props){
		props = props !== undefined ? props : {};

		this.Time = props.Time || 0;
		this.ThreadID = props.ThreadID || 0;
		this.StackID = props.StackID || 0;
	};

	ThreadWakeEvent.prototype = {
		Code: 0x05,
		read: function(stream){
			this.Time = stream.readInt();
			this.ThreadID = stream.readInt();
			this.StackID = stream.readInt();
		},
		write: function(stream){
			stream.writeInt(this.Time);
			stream.writeInt(this.ThreadID);
			stream.writeInt(this.StackID);
		}
	};

	Event.ThreadStop = ThreadStopEvent;
	ThreadStopEvent.Code = 0x06;
	EventCode.ThreadStop = 0x06;
	EventByCode[0x06] = ThreadStopEvent;
	function ThreadStopEvent(props){
		props = props !== undefined ? props : {};

		this.Time = props.Time || 0;
		this.ThreadID = props.ThreadID || 0;
		this.StackID = props.StackID || 0;
	};

	ThreadStopEvent.prototype = {
		Code: 0x06,
		read: function(stream){
			this.Time = stream.readInt();
			this.ThreadID = stream.readInt();
			this.StackID = stream.readInt();
		},
		write: function(stream){
			stream.writeInt(this.Time);
			stream.writeInt(this.ThreadID);
			stream.writeInt(this.StackID);
		}
	};

	Event.Begin = BeginEvent;
	BeginEvent.Code = 0x07;
	EventCode.Begin = 0x07;
	EventByCode[0x07] = BeginEvent;
	function BeginEvent(props){
		props = props !== undefined ? props : {};

		this.Time = props.Time || 0;
		this.ThreadID = props.ThreadID || 0;
		this.StackID = props.StackID || 0;
		this.ID = props.ID || 0;
	};

	BeginEvent.prototype = {
		Code: 0x07,
		read: function(stream){
			this.Time = stream.readInt();
			this.ThreadID = stream.readInt();
			this.StackID = stream.readInt();
			this.ID = stream.readInt();
		},
		write: function(stream){
			stream.writeInt(this.Time);
			stream.writeInt(this.ThreadID);
			stream.writeInt(this.StackID);
			stream.writeInt(this.ID);
		}
	};

	Event.End = EndEvent;
	EndEvent.Code = 0x08;
	EventCode.End = 0x08;
	EventByCode[0x08] = EndEvent;
	function EndEvent(props){
		props = props !== undefined ? props : {};

		this.Time = props.Time || 0;
		this.ThreadID = props.ThreadID || 0;
		this.StackID = props.StackID || 0;
		this.ID = props.ID || 0;
	};

	EndEvent.prototype = {
		Code: 0x08,
		read: function(stream){
			this.Time = stream.readInt();
			this.ThreadID = stream.readInt();
			this.StackID = stream.readInt();
			this.ID = stream.readInt();
		},
		write: function(stream){
			stream.writeInt(this.Time);
			stream.writeInt(this.ThreadID);
			stream.writeInt(this.StackID);
			stream.writeInt(this.ID);
		}
	};

	Event.Start = StartEvent;
	StartEvent.Code = 0x09;
	EventCode.Start = 0x09;
	EventByCode[0x09] = StartEvent;
	function StartEvent(props){
		props = props !== undefined ? props : {};

		this.Time = props.Time || 0;
		this.ThreadID = props.ThreadID || 0;
		this.StackID = props.StackID || 0;
		this.ID = props.ID || 0;
	};

	StartEvent.prototype = {
		Code: 0x09,
		read: function(stream){
			this.Time = stream.readInt();
			this.ThreadID = stream.readInt();
			this.StackID = stream.readInt();
			this.ID = stream.readInt();
		},
		write: function(stream){
			stream.writeInt(this.Time);
			stream.writeInt(this.ThreadID);
			stream.writeInt(this.StackID);
			stream.writeInt(this.ID);
		}
	};

	Event.Finish = FinishEvent;
	FinishEvent.Code = 0x0A;
	EventCode.Finish = 0x0A;
	EventByCode[0x0A] = FinishEvent;
	function FinishEvent(props){
		props = props !== undefined ? props : {};

		this.Time = props.Time || 0;
		this.ThreadID = props.ThreadID || 0;
		this.StackID = props.StackID || 0;
		this.ID = props.ID || 0;
	};

	FinishEvent.prototype = {
		Code: 0x0A,
		read: function(stream){
			this.Time = stream.readInt();
			this.ThreadID = stream.readInt();
			this.StackID = stream.readInt();
			this.ID = stream.readInt();
		},
		write: function(stream){
			stream.writeInt(this.Time);
			stream.writeInt(this.ThreadID);
			stream.writeInt(this.StackID);
			stream.writeInt(this.ID);
		}
	};

	Event.Sample = SampleEvent;
	SampleEvent.Code = 0x0B;
	EventCode.Sample = 0x0B;
	EventByCode[0x0B] = SampleEvent;
	function SampleEvent(props){
		props = props !== undefined ? props : {};

		this.Time = props.Time || 0;
		this.ThreadID = props.ThreadID || 0;
		this.StackID = props.StackID || 0;
		this.Values = props.Values || new Array();
	};

	SampleEvent.prototype = {
		Code: 0x0B,
		read: function(stream){
			this.Time = stream.readInt();
			this.ThreadID = stream.readInt();
			this.StackID = stream.readInt();
			this.Values = stream.readValues();
		},
		write: function(stream){
			stream.writeInt(this.Time);
			stream.writeInt(this.ThreadID);
			stream.writeInt(this.StackID);
			stream.writeValues(this.Values);
		}
	};

	Event.Snapshot = SnapshotEvent;
	SnapshotEvent.Code = 0x0C;
	EventCode.Snapshot = 0x0C;
	EventByCode[0x0C] = SnapshotEvent;
	function SnapshotEvent(props){
		props = props !== undefined ? props : {};

		this.Time = props.Time || 0;
		this.ThreadID = props.ThreadID || 0;
		this.StackID = props.StackID || 0;
		this.ID = props.ID || 0;
		this.ContentKind = props.ContentKind || 0;
		this.Content = props.Content || new Uint8Array();
	};

	SnapshotEvent.prototype = {
		Code: 0x0C,
		read: function(stream){
			this.Time = stream.readInt();
			this.ThreadID = stream.readInt();
			this.StackID = stream.readInt();
			this.ID = stream.readInt();
			this.ContentKind = stream.readByte();
			this.Content = stream.readBlob();
		},
		write: function(stream){
			stream.writeInt(this.Time);
			stream.writeInt(this.ThreadID);
			stream.writeInt(this.StackID);
			stream.writeInt(this.ID);
			stream.writeByte(this.ContentKind);
			stream.writeBlob(this.Content);
		}
	};

	Event.Info = InfoEvent;
	InfoEvent.Code = 0x0D;
	EventCode.Info = 0x0D;
	EventByCode[0x0D] = InfoEvent;
	function InfoEvent(props){
		props = props !== undefined ? props : {};

		this.ID = props.ID || 0;
		this.Name = props.Name || '';
		this.ContentKind = props.ContentKind || 0;
		this.Content = props.Content || new Uint8Array();
	};

	InfoEvent.prototype = {
		Code: 0x0D,
		read: function(stream){
			this.ID = stream.readInt();
			this.Name = stream.readUTF8();
			this.ContentKind = stream.readByte();
			this.Content = stream.readBlob();
		},
		write: function(stream){
			stream.writeInt(this.ID);
			stream.writeUTF8(this.Name);
			stream.writeByte(this.ContentKind);
			stream.writeBlob(this.Content);
		}
	};


	var ContentKind = {
		Invalid: 0x00,
		Thread: 0x01,
		Stack: 0x02,
		Text: 0x10,
		JSON: 0x11,
		BLOB: 0x12,
		Image: 0x13,
		User: 0x20,
	};

	return {
		Version: 1,
		Event: Event,
		EventCode: EventCode,
		EventByCode: EventByCode,
		ContentKind: ContentKind
	};
})