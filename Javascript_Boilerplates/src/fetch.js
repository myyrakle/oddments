async function fetchGet(url, option={})
{
      const response = await fetch(url, {
          method:'GET',
          headers: option.headers
      });

      return response;
}


async function fetchPost(url, option={})
{
      const response = await fetch(url, {
          method:'POST',
          headers: option.headers,
          body: JSON.stringify(option.body)
      });

      return response;
}
