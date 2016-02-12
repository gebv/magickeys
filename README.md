# magic keys

[![Build Status](https://travis-ci.org/gebv/magickeys.svg?branch=master)](https://travis-ci.org/gebv/magickeys)

Эксперементальный проект.
Основная идея - в качестве ключа (в таблице где хранятся данные) используется массив ключей. Для хранения данных используется РСУБД Postgres. 

## api

| URL | Описание |
| ---| --- |
| /api/v1/values/ | CRUD |
| /api/v1/values/search/eq/{keys} | Поиск записей по точному совпадению ключей |
| /api/v1/values/search/any/{keys} | Поиск всех записей в которых встрачется хотя бы один ключ из keys |
| /api/v1/values/search/contains/{keys} | Поиск всех записей в которых keys является подмножеством |

Про any, contains см. подробней в описании [postgresql array functions](http://www.postgresql.org/docs/9.4/static/functions-array.html) для операторов **&&** и **@>**.

# Примеры использования

В качестве фронтенда используется [mithril](http://mithril.js.org)

* многоуровневый TODO лист
* Таблицы (с конструктором)

### многоуровневый TODO лист

[todo list](web/example/todolist.html)

![многоуровневый todo list](images/magickey_todolist.gif)

```
magickeys=# SELECT keys, value, props FROM values;
-[ RECORD 1 ]------------------------------------------------------------------------------------------
keys  | {examples,simplelist}
value | 1
props | {"ts": "ts1455291704475", "done": "no"}
-[ RECORD 2 ]------------------------------------------------------------------------------------------
keys  | {examples,simplelist}
value | 2
props | {"ts": "ts1455291704949", "done": "no"}
-[ RECORD 3 ]------------------------------------------------------------------------------------------
keys  | {examples,simplelist}
value | 3
props | {"ts": "ts1455291705508", "done": "no"}
-[ RECORD 4 ]------------------------------------------------------------------------------------------
keys  | {examples,simplelist,1f0803ee-d19f-11e5-b356-10ddb19b9d24}
value | 3.1
props | {"ts": "ts1455291708861", "done": "no"}
-[ RECORD 5 ]------------------------------------------------------------------------------------------
keys  | {examples,simplelist,1f0803ee-d19f-11e5-b356-10ddb19b9d24,22b85341-d19f-11e5-b356-10ddb19b9d24}
value | 3.3.1
props | {"ts": "ts1455291714655", "done": "no"}
-[ RECORD 6 ]------------------------------------------------------------------------------------------
keys  | {examples,simplelist,1e6bf1b0-d19f-11e5-b356-10ddb19b9d24}
value | 1.1
props | {"ts": "ts1455291717220", "done": "no"}
-[ RECORD 7 ]------------------------------------------------------------------------------------------
keys  | {examples,simplelist,1f0803ee-d19f-11e5-b356-10ddb19b9d24}
value | 3.2
props | {"ts": "ts1455291709930", "done": "yes"}
-[ RECORD 8 ]------------------------------------------------------------------------------------------
keys  | {examples,simplelist,1f0803ee-d19f-11e5-b356-10ddb19b9d24}
value | 3.3
props | {"ts": "ts1455291711696", "done": "yes"}

magickeys=#
```

### Таблицы (с конструктором)

Динамические поля\столбцы управляющийся через конструктор.

## database schema

``` sql
create function sort_text_array(text[]) returns text[][] as $$
    select array_agg(n) from (select n from unnest($1) as t(n) order by n) as a;
$$ language sql immutable;

CREATE TABLE values (
    value_id uuid NOT NULL PRIMARY KEY,
    keys text[],
    value text,
    props jsonb NOT NULL DEFAULT '{}', -- Вспомогательное поле для хранения расширенных значений
    flags text[], -- Вспомогательное поле для хранение расширенных значений

    is_enabled boolean DEFAULT true,
    is_removed boolean DEFAULT false,
    created_at timestamp,
    updated_at timestamp DEFAULT now()
);
CREATE INDEX values_keys_idx on values USING GIN (keys);
CREATE UNIQUE INDEX values_keys_ifuniq_idx on values (sort_text_array(keys))
    WHERE keys @> '{uniq}';

```
