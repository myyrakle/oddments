function isString(value)
{
  return value === String(value);
}

function isNumber(value)
{
  return value === Number(value);
}

function isInteger(value)
{
  return (value === Number(value)) && value%1 === 0;
}

function isFloat(value)
{
  return (value === Number(value)) && value%1 !== 0;
}

function isBoolean(value)
{
  return value === Boolean(value);
}

function isNull(value)
{
  return value === null;
}

function isUndefined(value)
{
  return value === undefined;
}

function isArray(value)
{
  Array.isArray(value);
}

function isObject(value)
{
  return value === Object(value);
}

function isFunction(value)
{
  return value instanceof Function;
}
