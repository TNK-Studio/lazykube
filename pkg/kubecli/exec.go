package kubecli

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/exec"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"time"

	"k8s.io/kubectl/pkg/util/templates"
)

// Note: Copy code because of argsLenAtDash always "-1"

var (
	execExample = templates.Examples(i18n.T(`
		# Get output from running 'date' command from pod mypod, using the first container by default
		kubectl exec mypod -- date

		# Get output from running 'date' command in ruby-container from pod mypod
		kubectl exec mypod -c ruby-container -- date

		# Switch to raw terminal mode, sends stdin to 'bash' in ruby-container from pod mypod
		# and sends stdout/stderr from 'bash' back to the client
		kubectl exec mypod -c ruby-container -i -t -- bash -il

		# List contents of /usr from the first container of pod mypod and sort by modification time.
		# If the command you want to execute in the pod has any flags in common (e.g. -i),
		# you must use two dashes (--) to separate your command's flags/arguments.
		# Also note, do not surround your command and its flags/arguments with quotes
		# unless that is how you would execute it normally (i.e., do ls -t /usr, not "ls -t /usr").
		kubectl exec mypod -i -t -- ls -t /usr

		# Get output from running 'date' command from the first pod of the deployment mydeployment, using the first container by default
		kubectl exec deploy/mydeployment -- date

		# Get output from running 'date' command from the first pod of the service myservice, using the first container by default
		kubectl exec svc/myservice -- date
		`))
)

const (
	defaultPodExecTimeout = 60 * time.Second
)

func NewCmdExec(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	options := &exec.ExecOptions{
		StreamOptions: exec.StreamOptions{
			IOStreams: streams,
		},

		Executor: &exec.DefaultRemoteExecutor{},
	}
	cmd := &cobra.Command{
		Use:                   "exec (POD | TYPE/NAME) [-c CONTAINER] [flags] -- COMMAND [args...]",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Execute a command in a container"),
		Long:                  "Execute a command in a container.",
		Example:               execExample,
		Run: func(cmd *cobra.Command, args []string) {
			// Note: Copy code because of argsLenAtDash always "-1"
			// argsLenAtDash := cmd.ArgsLenAtDash()
			argsLenAtDash := 1
			cmdutil.CheckErr(options.Complete(f, cmd, args, argsLenAtDash))
			cmdutil.CheckErr(options.Validate())
			cmdutil.CheckErr(options.Run())
		},
	}
	cmdutil.AddPodRunningTimeoutFlag(cmd, defaultPodExecTimeout)
	cmdutil.AddJsonFilenameFlag(cmd.Flags(), &options.FilenameOptions.Filenames, "to use to exec into the resource")
	// TODO support UID
	cmd.Flags().StringVarP(&options.ContainerName, "container", "c", options.ContainerName, "Container name. If omitted, the first container in the pod will be chosen")
	cmd.Flags().BoolVarP(&options.Stdin, "stdin", "i", options.Stdin, "Pass stdin to the container")
	cmd.Flags().BoolVarP(&options.TTY, "tty", "t", options.TTY, "Stdin is a TTY")
	return cmd
}

// Exec Exec
func (cli *KubeCLI) Exec(streams genericclioptions.IOStreams, args ...string) *Cmd {
	cmd := NewCmdExec(cli.factory, streams)
	return NewCmd(cmd, args, streams)
}
