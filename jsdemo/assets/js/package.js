var global = window || this;
global.package = function package(name, setup){
	if(name == ""){
		throw new Error("package name cannot be empty");
	}

	var info = package.find(name);
	var inject = setup(info.namespace);
	for(var name in inject){
		info.namespace[name] = inject[name];
	}
};

global.package.find = function(name){
	var created = false;
	var path = name.split(".");
	var namespace = global;
	for(var i = 0; i < path.length; i++){
		var token = path[i];
		var next = namespace[token];
		if(next){
			created = false;
		} else {
			next = {};
			namespace[token] = next;
			created = true;
		}
		namespace = next;
	}

	return {
		namespace: namespace,
		created: created
	};
}

global.depends = function depends(name){
	var info = package.find(name);
	if(info.created){
		throw new Error("package " + name + " not loaded.");
	}
	return info.namespace;
}