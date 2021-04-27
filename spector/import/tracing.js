package("spector.import.tracing", function(){
	function Stream(data) {
		this.stage = spector.Stream.Stage.Reading;
	}

	Stream.prototype = {
		next: function(){
			if(this.buffer_.length <= 0) {
				return undefined;
			}
			return this.buffer_.shift();
		}
	};

	return {
		Stream: Stream
	};
});