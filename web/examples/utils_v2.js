// Helpful functions
	var api = "/api/v1"
	var buildUrlFindByKeysFromMode = function(mode, keys) {
		// mode {contains, eq, any}

		return api+"/values/search/"+encodeURIComponent(mode)+"/"+encodeURIComponent(keys.join(","));
	};

	var keysToString = function(keys) {
		return _.sortBy(keys).join(":")
	}


	var buildUrlItem = function(value_id, fields) {
		return api+"/values/"+encodeURIComponent(value_id)+
			(fields && fields.length > 0? 
				"?"+m.route.buildQueryString({fields: fields}):
				"");
	}

	var humanTimeToTimestamp = function(dateString) {
		return Math.round(Date.parse(dateString) / 1000)
	}

	var Manager = function(config) {
		this.isLoading = m.prop(false);
		this.items = m.prop([]);
		this.keys = [];
		this.components = {};

		_.merge(this, config);

		console.log(this);
		// TODO: для каждой комбинации ключей свой дозагрузчик (сохранение курсора)
	}

	Manager.prototype.makeUpdate = function(value_id, item) {
		var _item = _.find(this.items(), {value_id: value_id});

		if (!_item) {
			var res = {message: "Ошибка поиска"};
			return {
				then: function(cbSuccess) {
					cbSuccess(res);
				}
			}
		}

		var _newDataItem = {
			value: _item.data.value(),
			keys: _item.keys,
		}

		_newDataItem = _.assignIn({}, _newDataItem, item);

		_item.isLoading(true);
		m.endComputation();

		return m.request({url: buildUrlItem(value_id, ["value", "keys"]), "method": "PUT", data: _newDataItem})
			.then(
				function(res){
					_item.isLoading(false);

					if (res.status_code == 200 && res.data && res.data.value_id && res.data.value_id.length != 0) {
						_item.data.value(res.data.value);
						_item._value = res.data.value; // for hasNotChanged
						_item.data.updated_at(res.data.updated_at);
					} else {
						var message = ["Ошибка обновления"];
						message.push("value_id: "+value_id)
						message.push("item: "+JSON.stringify(item))
						message.push("response: "+JSON.stringify(res))
						alert(message.join("\n"));
					}

					return _item;
				}, 
				function(res){
					var message = ["Ошибка обновления"];
					message.push("value_id: "+value_id)
					message.push("item: "+JSON.stringify(item))
					message.push("response: "+JSON.stringify(res))
					alert(message.join("\n"));
				}.bind(this)
			)
			.catch(function(e){
				console.error("makeUpdate", e)
			})
	}

	Manager.prototype.makeDelete = function(value_id) {
		var _item = _.find(this.items(), {value_id: value_id});

		if (!_item) {
			var res = {status_code: 400, message: "Ошибка поиска"};
			return {
				then: function(cbSuccess) {
					cbSuccess(res);
				}
			}
		}

		if (_item.isNew()) {
			_item.data.is_removed(true);
			
			var res = {status_code: 200, message: "Успешно удалено (временная запись на клиенте)"};
			return {
				
				then: function(cbSuccess) {
					cbSuccess(res);
				}
			}
		}

		_item.isLoading(true);

		return m.request({url: buildUrlItem(value_id), "method": "DELETE"})
			.then(function(res){
				_item.isLoading(false);

				if (res.status_code == 200) {
					// TODO: code
					console.log("makeDelete", res)
					_item.data.is_removed(true);
				}

			}, function(res){

			})
			.catch(function(e){
				console.error("makeDelete", e)
			})
	}

	Manager.prototype.getComponent = function(name, config) {
		var config = _.merge({}, this, config)
		return m.component(this.components[name], config)
	}

	Manager.prototype.getByKeys = function(keys) {

		return _(this.items())
			.filter(function(_item){
				return keysToString(_item.keys) == keysToString(keys) && _item.data.is_removed() == false;
			})
			.orderBy(function(_item){return humanTimeToTimestamp(_item.data.created_at())}, "desc")
			.value();
	},

	Manager.prototype.getByIds = function(ids) {
		return _(this.items())
			.filter(function(_item){
				return _.intersection(ids, [_item.value_id]).length == 1 && _item.data.is_removed() == false;
			})
			.orderBy(function(_item){return humanTimeToTimestamp(_item.data.created_at())}, "desc")
			.value();
	}

	Manager.prototype.makeCreate = function(item) {
		console.log("makeCreate", item)
		var _item = this.add(item);

		_item.isLoading(true);
		
		var _newDataItem = {
			value: _item.data.value(),
			keys: _item.keys
		}

		return m.request({url: buildUrlItem("", []), "method": "POST", data: _newDataItem})
			.then(
				function(res) {
					_item.isLoading(false);

					if (res.status_code == 200 && res.data && res.data.value_id && res.data.value_id.length != 0) {
						_item.value_id = res.data.value_id;

						_item.data.value_id(res.data.value_id);
						_item._value = res.data.value; // for hasNotChanged

						_item.data.created_at(res.data.created_at);
						_item.data.updated_at(res.data.updated_at);

						_item.isNew(false);
						console.log("makeCreate:rsponse", _item)
					} else {
						var message = ["Ошибка создания"];
						message.push("item: "+JSON.stringify(item))
						message.push("response: "+JSON.stringify(res))
						alert(message.join("\n"));
					}

					return _item;
				}.bind(this), 
				function(res) {
					var message = ["Ошибка создания"];
					message.push("item: "+JSON.stringify(item))
					message.push("response: "+JSON.stringify(res))
					alert(message.join("\n"));
				}
			)
			.catch(function(e) {
				console.error("makeCreate", e)
			})
	}

	Manager.prototype.new = function(item) {
		item = _.assign({}, item, {value_id: "new:ts:"+_.now()});

		_item = new Item(item, this);
		_item.isNew(true)

		this.items().push(_item);

		return _item
	}

	Manager.prototype.add = function(item) {
		_item = _.find(this.items(), {value_id: item.value_id});

		if (!_item) {
			_item = new Item(item, this);	
			this.items().push(_item);
		}

		return _item;
	}

	Manager.prototype.loadByKeys = function(keys) {
		// contains

		this.isLoading(true);

		return m.request({url: buildUrlFindByKeysFromMode("contains", keys || this.keys)})
			.then(function(res){
				this.isLoading(false);

				if (res.status_code == 200 && res.data) {

					_.each(res.data, this.add.bind(this));
				} else {
					var message = ["Ошибка создания"];
					message.push("item: "+JSON.stringify(item))
					message.push("response: "+JSON.stringify(res))
					alert(message.join("\n"));
				}

				return this
			}.bind(this), function(res){
				var message = ["Ошибка создания"];
				message.push("item: "+JSON.stringify(item))
				message.push("response: "+JSON.stringify(res))
				alert(message.join("\n"));
			})
			.catch(function(e) {
				console.error("loadByKeys", e);
			}.bind(this))
	}

	var NewItem = function(item){
		return {
			value_id: m.prop(item && item.value_id || ""),
			value: m.prop(item && item.value || {}),
			keys: m.prop(item && item.keys || []),
			created_at: m.prop(item && item.created_at || (new Date()).toString()),
			updated_at: m.prop(item && item.updated_at || (new Date()).toString()),
			is_removed: m.prop(item && item.is_removed || false),
		}
	}

	var Item = function(item, manager) {
		this.data = NewItem(item);
		this.value_id = this.data.value_id();
		this.keys = this.data.keys();

		this._value = _.extend({}, this.data.value());
		this._props = m.prop({});

		this.isLoading = m.prop(false);
		this.isNew = m.prop(false);

		this.manager = manager;

	}

	Item.prototype.hasNotChanged = function() {
		var f1 = _.sortBy(_.reduce(this.data.value(), function(res, v, i){res.push({v: v, i: i}); return res;}, []), "i").map(function(value){return value.v+":"+value.i}).join(",")
		var f2 = _.sortBy(_.reduce(this._value, function(res, v, i){res.push({v: v, i: i}); return res;}, []), "i").map(function(value){return value.v+":"+value.i}).join(",")

		return f1 == f2
	}

	// Создание нового элемента
	Item.prototype.create = function() {
		return this.manager.makeCreate({value_id: this.value_id, value: this.data.value(), keys: this.data.keys})
	}

	Item.prototype.update = function() {
		return this.manager.makeUpdate(this.value_id, {value: this.data.value(), keys: this.data.keys})
	}

	Item.prototype.delete = function(config) {
		var shortTitle = _.truncate(config.confirm_title || this.value_id, {'length': 24});
		if (!confirm("Подтвердите удаление \""+shortTitle+"\"?")) {
			return
		}
		return this.manager.makeDelete(this.value_id)
	}