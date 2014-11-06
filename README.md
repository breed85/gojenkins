gojenkins
=========

Go program for launching a Jenkins slave. Supports JNLP connections or Swarm plugin connections. 


Usage
=========
usage: gojenkins

   This command will launch a jenkins slave. Command line options override the environment.
   All command line options can be set via environment variables prefixed with 'SLAVE_'.

  -executors=2:
	Number of executors to use for the node. Requires -swarm

  -file="":
	Config file to use. The content should be json. The file will be automatically monitored for changes.
	Any settings in the file will take precedence over command line flags.

  -home="~/go/src/github.com/breed85/gojenkins":
	Jenkins working directory.

  -labels="":
	Labels to apply to the node. Requires -swarm. Can be a space separated list.

  -lock=false:
	Create a lock file during execution in -home directory with name [name].lock

  -log="":
	Log file to output to. Default is STDERR.

  -mode="normal":
	Mode to set for the slave node. Valid values are 'normal' (utilize the slave as much as possible)
	or 'exclusive' (leave this machine for tied jobs only). Requires -swarm

  -name="localhost":
	Name of the host on Jenkins. When used with -swarm, the name will be used to create a node.
	Otherwise, the node [name] must exist on the master already.

  -password="":
	Password to log into the jenkins system. Requires -swarm

  -server="":
	Jenkins server to use. Ex. http://localhost:8080

  -swarm=false:
	Use swarm client to connect to jenkins.

  -swarmversion="1.15":
	Version of swarm client to use. Requires -swarm

  -username="":
	Username to log into the jenkins system. Requires -swarm
