# Zocker
A simple docker implementation in Go for learning purposes.

## Follow along

### Step 1
First, start a Go project. Create a go mod file, a cmd folder which will contain the zocker.go file. 

To zocker.go, add the most basic imports - fmt & OS, since we will be using all sorts of OS functionalities. 

I've also added _ import, because it is a really cool way to make cli apps, and it's easy to use as well.

### Step 2

Now, let's write the most basic code of zocker. 

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
Increasing the isolation of the containers, we've now moved on to isolating serveral other namespaces - Process ID, Mount, and IPC (inter-process communication),

	`process.SysProcAttr = &syscall.SysProcAttr{
							Cloneflags: syscall.CLONE_NEWUTS |
							syscall.CLONE_NEWPID |
							syscall.CLONE_NEWNS |
							syscall.CLONE_NEWIPC,
						}`

We've seperated them by using a bitwise OR. You can check the isolation by running `ps aux` for PID.

### Step 6

It's time for filesystel isolation.

Docker uses a sort of virtual filesystem, where multiple files are stored in a single location, overlayed on each other. In this manner, the container only takes up on location, and all its contents are stored in it. 

In linux, we can achieve this by using `chroot` tool. 

We are going to create a snapshot of the Btrfs volume (Btrfs (B-tree file system) is a file system for Linux that's designed for large storage systems), and this will act is the container's isolated root file system.

Note that we have added a new import.