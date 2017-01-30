#include <windows.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <objbase.h>
#include <shlobj.h>

int CreateShortcut(char *shortcutA, char *path, char *args) {
	IShellLink*   pISL;
	IPersistFile* pIPF;
	HRESULT       hr;

	hr = CoInitializeEx(NULL, COINIT_MULTITHREADED);
	if (!SUCCEEDED(hr)) {
		return FALSE;
	}

	// Shortcut filename: convert ANSI to unicode
	WORD shortcutW[MAX_PATH];
	int nChar = MultiByteToWideChar(CP_ACP, 0, shortcutA, -1, shortcutW, MAX_PATH);

	hr = CoCreateInstance(&CLSID_ShellLink, NULL, CLSCTX_INPROC_SERVER, &IID_IShellLink, (LPVOID*)&pISL);
	if (!SUCCEEDED(hr)) {
		return FALSE;
	}

	// See https://msdn.microsoft.com/en-us/library/windows/desktop/bb774950(v=vs.85).aspx
	hr = pISL->lpVtbl->SetPath(pISL, path);
	if (!SUCCEEDED(hr)) {
		return FALSE;
	}

	hr = pISL->lpVtbl->SetArguments(pISL, args);
	if (!SUCCEEDED(hr)) {
		return FALSE;
	}

	// Save the shortcut
	hr = pISL->lpVtbl->QueryInterface(pISL, &IID_IPersistFile, (void**)&pIPF);
	if (!SUCCEEDED(hr)) {
		return FALSE;
	}

	hr = pIPF->lpVtbl->Save(pIPF, shortcutW, FALSE);
	if (!SUCCEEDED(hr)) {
		return FALSE;
	}

	pIPF->lpVtbl->Release(pIPF);
	pISL->lpVtbl->Release(pISL);
	return TRUE;
}