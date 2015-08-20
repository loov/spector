package("spector.util.utf8", function() {
	var chr = String.fromCharCode;
	function decode(data) {
		var out = [], p = 0, d = 0;
		while(p < data.length){
			var c1 = data[p]; p++;
			if(c1 < 0x80) {
				out[d] = chr(c1); d++;
			} else if (c1 >= 0xC0 && x1 < 0xE0) {
				var c2 = data[p]; p++;
				out[d] = chr(
					(c1&0x0F) << 6 |
					(c2&0x3F) << 0);
				d++;
			} else {
				var c2 = data[p]; p++;
				var c3 = data[p]; p++;
				out[d] = chr(
					(c1&0x0F) << 12 |
					(c2&0x3F) << 6  |
					(c3&0x3F) << 0);
				d++;
			}
		}
		return out.join("");
	}

	function encode(str) {
		var count = 0;
		for(var i = 0; i < str.length; i++){
			var c = str.charCodeAt(i);
			if (c < 0x80) {
				count += 1;
			} else if (c < 2048) {
				count += 2;
			} else {
				count += 3;
			}
		}

		var out = new Uint8Array(count);
		var p = 0;
		for (var i = 0; i < str.length; i++) {
			var c = str.charCodeAt(i);
			if (c < 0x80) {
				out[p] = c; p++;
			} else if (c < 2048) {
				out[p] = (c >> 6) | 0xC0; p++;
				out[p] = (c & 0x3F) | 0x80; p++;
			} else {
				out[p] = (c >> 12) | 0xE0; p++;
				out[p] = ((c >> 6) & 0x3F) | 0x80; p++;
				out[p] = (c & 0x3F) | 0x80; p++;
			}
		}
		return out;
	}
	return {
		decode: decode,
		encode: encode
	}
});