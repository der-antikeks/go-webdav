# webdav

__TODO (litmus tests):__
* basic..... pass
* copymove.. pass
* props..... FAILED
	* init.................. pass
	* begin................. pass
	* propfind_invalid...... pass
	* propfind_invalid2..... pass
	* propfind_d0........... pass
	* propinit.............. pass
	* propset............... FAIL
	* propget............... SKIPPED
	* propextended.......... pass
	* propmove.............. SKIPPED
	* propget............... SKIPPED
	* propdeletes........... SKIPPED
	* propget............... SKIPPED
	* propreplace........... SKIPPED
	* propget............... SKIPPED
	* propnullns............ SKIPPED
	* propget............... SKIPPED
	* prophighunicode....... SKIPPED
	* propget............... SKIPPED
	* propremoveset......... SKIPPED
	* propget............... SKIPPED
	* propsetremove......... SKIPPED
	* propget............... SKIPPED
	* propvalnspace......... SKIPPED
	* propwformed........... pass
	* propinit.............. pass
	* propmanyns............ FAIL
	* propget............... FAIL
	* propcleanup........... pass
	* finish................ pass
* locks..... SKIPPED
* http...... SKIPPED

__future:__
* server
* client
* replace litmus test with plain go tests
