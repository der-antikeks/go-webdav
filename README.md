# webdav

__TODO (litmus tests):__
* basic..... pass
	* options............... pass
	* put_get............... pass
	* put_get_utf8_segment.. pass
	* put_no_parent......... pass
	* mkcol_over_plain...... pass
	* delete................ pass
	* delete_null........... pass
	* delete_fragment....... pass
	* mkcol................. pass
	* mkcol_again........... pass
	* delete_coll........... pass
	* mkcol_no_parent....... pass
	* mkcol_with_body....... pass
* copymove.. FAIL
	* copy_init............. pass
	* copy_simple........... pass
	* copy_overwrite........ FAIL
	* copy_nodestcoll....... WARNING
	* copy_cleanup.......... pass
	* copy_coll............. FAIL
	* copy_shallow.......... FAIL
	* move.................. FAIL
	* move_coll............. FAIL
	* move_cleanup.......... pass
* props..... SKIPPED
* locks..... SKIPPED
* http...... SKIPPED

__future:__
* server
* client
* replace litmus test with plain go tests
