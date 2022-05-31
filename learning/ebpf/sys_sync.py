from bcc import BPF

prog = '''
  int kprobe__sys_sync(void *ctx) {
    bpf_trace_printk("sys_sync called \\n");
    return 0;
  }
'''

print('Tracing sys_sync... Ctrl-C to end')
BPF(text=prog).trace_print()
