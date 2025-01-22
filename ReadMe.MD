# Zocker
A simple docker implementation in Go for learning purposes.

## Follow along

### Step 1

### Step 2

### Step 3
Next, we've added a Create command for creating a container. Giving the container a name is mandatory, otherwise our code returns an error. We;ve defined the directory for the container, being "/tmp/%s". Using the OS package, we create that directory and give it the permission 0755.

`err := os.Mkdir(containerDir, 0755)`

0755 indicates Read, write and execute permission for the owner.

Like always, we've added error handling, by including possible errors at various stages here. 

If all goes well, a confirmation message is printed.

#### Current Progress
 `chhax1618@apprentice:~/code/learn/docker/zocker$ go run cmd/zocker.go create name
Container 'name' created at /tmp/name`

### Step 4
Now, it's time for linux namespace. Docker containers have a concept calle dprocess isolation. That is, the processes going on in a docker container are isolated from the host. They can't access resources outside the container environment unless specified.

Docker containers isolate the folloring namespaces - UTS, User, File System, PID, Mount, network, etc.

Here, we have implemented the UTS isolation.

                    process.SysProcAttr = &syscall.SysProcAttr{
						Cloneflags: syscall.CLONE_NEWUTS,
					}

This piece of code isolates the hostname and similar identifiers. The tag CLONE_NEWUTS is used for this purpose. The system call for the system process attributes (SysProcAttr) allows the specification of low level, OS specific attributes for the process.

Now when we use the run command, the host is isolated from the main host (our PC).


### Step 5
