#!/usr/bin/env python3
"""
Simple development script for mArchive
Runs both frontend and backend with hot reload
"""

import os
import sys
import subprocess
import signal
import threading
import atexit
from pathlib import Path


# Global process list for cleanup
processes = []


def start_process(cmd, cwd, name):
    """Start a process as a true child process."""
    try:
        if os.name == 'nt':
            # Windows: Use job objects to ensure child processes die with parent
            import ctypes
            from ctypes import wintypes
            
            proc = subprocess.Popen(
                cmd, cwd=cwd,
                stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                text=True, bufsize=1, shell=True,
                creationflags=subprocess.CREATE_NEW_PROCESS_GROUP
            )
            
            # Create job object and assign process to it
            job = ctypes.windll.kernel32.CreateJobObjectW(None, None)
            if job:
                # Configure job to kill all processes when parent dies
                job_info = ctypes.create_string_buffer(16)  # JOBOBJECT_EXTENDED_LIMIT_INFORMATION size
                ctypes.windll.kernel32.SetInformationJobObject(
                    job, 9, job_info, 16  # JobObjectExtendedLimitInformation = 9
                )
                
                # Assign process to job
                ctypes.windll.kernel32.AssignProcessToJobObject(job, proc._handle)
                
        else:
            # Unix: Set parent death signal
            def set_pdeathsig():
                import ctypes
                import ctypes.util
                libc = ctypes.CDLL(ctypes.util.find_library("c"))
                PR_SET_PDEATHSIG = 1
                libc.prctl(PR_SET_PDEATHSIG, signal.SIGTERM)
                
            proc = subprocess.Popen(
                cmd, cwd=cwd,
                stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                text=True, bufsize=1,
                preexec_fn=set_pdeathsig
            )
        
        print(f"🚀 {name} started (PID: {proc.pid})")
        processes.append(proc)
        return proc
        
    except Exception as e:
        print(f"❌ Failed to start {name}: {e}")
        return None


def monitor_output(proc, name):
    """Monitor process output in a separate thread."""
    if not proc:
        return
    
    def _monitor():
        try:
            for line in iter(proc.stdout.readline, ''):
                if line.strip():
                    print(f"[{name.upper()}] {line.rstrip()}")
        except:
            pass  # Process died, that's ok
    
    thread = threading.Thread(target=_monitor, daemon=True)
    thread.start()
    return thread


def cleanup_all():
    """Emergency cleanup function."""
    for proc in processes:
        try:
            if proc.poll() is None:
                proc.terminate()
                try:
                    proc.wait(timeout=2)
                except subprocess.TimeoutExpired:
                    proc.kill()
        except:
            pass


def main():
    project_root = Path(__file__).parent
    frontend_dir = project_root / "frontend"
    
    # Register cleanup function for emergency exits
    atexit.register(cleanup_all)
    
    # Set basic environment
    os.environ.update({'DEV_MODE': 'true', 'NODE_ENV': 'development'})
    
    print("🎯 Starting mArchive development environment...")
    
    # Start processes
    active_processes = {}
    
    # Start backend
    backend = start_process(
        ["go", "run", "-tags=debug", "cmd/server/main.go"], 
        project_root, "backend"
    )
    if backend:
        active_processes['backend'] = backend
    
    # Start frontend
    if frontend_dir.exists():
        frontend = start_process(
            ["npm", "run", "dev"], 
            frontend_dir, "frontend"
        )
        if frontend:
            active_processes['frontend'] = frontend
    else:
        print("⚠️  Frontend directory not found, skipping")
    
    if not active_processes:
        print("❌ No processes started successfully")
        sys.exit(1)
    
    # Monitor output
    for name, proc in active_processes.items():
        monitor_output(proc, name)
    
    print("🎉 Development environment running!")
    print("📱 Frontend: http://localhost:5173")
    print("🚀 Backend:  http://localhost:9520")
    print("Press Ctrl+C to stop")
    
    # Handle shutdown signals
    def signal_handler(signum, frame):
        print(f"\n🛑 Shutting down...")
        cleanup_all()
        sys.exit(0)
    
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)
    
    # Wait for processes
    try:
        while any(proc.poll() is None for proc in active_processes.values()):
            threading.Event().wait(1)
        print("❌ A process stopped unexpectedly")
    except KeyboardInterrupt:
        print("\n🛑 Interrupted by user")
    
    cleanup_all()
    print("👋 Development environment stopped")


if __name__ == "__main__":
    main()