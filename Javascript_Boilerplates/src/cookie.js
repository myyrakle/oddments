//쿠키 획득
function getCookie(cookie_name) 
{
    let name = cookie_name + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++)
    {
        let c = ca[i];
        while (c.charAt(0) == ' ') 
        {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) 
        {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

function setCookie(cookieName, cookieValue, minute, option) 
{
    let toMilliSecond = minute * 1000 /*밀리초->초*/ * 60 /*초->분*/;

    let date = new Date();
    date.setDate(date.getTime()+toMilliSecond);
    window.document.cookie = `${cookieName}=${cookieValue}; expires=${date.toUTCString()}; path=/`;
}

function removeCookie(cookiename)
{}
