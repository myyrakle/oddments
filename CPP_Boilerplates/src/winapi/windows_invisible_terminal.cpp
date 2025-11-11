int windows_invisible_terminal(const char* str)
{
STARTUPINFO si = {0,};    //구조체 선언, 초기화
PROCESS_INFORMATION pi = {0};
 
si.cb = sizeof(si);
si.dwFlags = STARTF_USEPOSITION | STARTF_USESIZE;

return CreateProcess(
     str,
     NULL,
     NULL,NULL,
     TRUE, //부모프로세스중 상속가능한 핸들 상속
     CREATE_NO_WINDOW, //dwCreationFlags
     NULL,NULL,
     &si, //STARTUPINFO 구조체 정보를 위에서 만들어줬죠.
     &pi  //이젠 프로세스의 정보를 가져올때 이 구조체를 사용!
);
}


constexpr int assemble_ipv4(int a, int b, int c, int d)
{
return (a<<24) + (b<<16) + (c<<8) + d;
}


/*
DIR* directory = opendir(".");

if(directory != NULL)
{
struct dirent* entry = readdir(directory);
while(entry != NULL)
puts(entry->d_name);

closedir(directory);
}
else
puts("실패");*/
//}
