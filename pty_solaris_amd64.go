package pty

// go types that mimic the natural C structure layout when compiling for 64 bit
// but without using cgo because we want to cross compile this easily

/*
struct strioctl {
        int     ic_cmd;                 // command
        int     ic_timout;              // timeout value
        int     ic_len;                 // length of data
        char    *ic_dp;                 // pointer to data
};
*/
type strioctl struct {
	ic_cmd    int32
	ic_timout int32
	ic_len    int32
	_         int32 //padding
	ic_dp     uintptr
}

/*
typedef struct pt_own {
        uid_t   pto_ruid;
        gid_t   pto_rgid;
} pt_own_t;
*/
type uid_t int32
type gid_t int32

type pt_own struct {
	pto_ruid uid_t
	pto_rgid gid_t
}
