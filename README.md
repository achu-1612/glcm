# glcm (go-routine life cycle management)
A way to manage complete life cycle of go-routines

TODO:
- support for timeout for go-routine shutdowns (if possible)
- do i really need a wg in the base runner?
- proper error handling for the pre and post hooks for service
- service dependency?
- unified way for stop and stopWait for service wrapper
- auto restart config for services
- exponential back-off restart for restart of service
