# Usage

## Execution

run:
```
make run
```

Also parameters host and port for the queue can also be changed (hereinafter **run-options**):
<pre><code><i>*only host*</i>
make run host=somehost
<i>*only port*</i>
make run port=123
<i>*both*</i>
make run host=somehost port=123
</code></pre>

build:
```
make build
```

build and run:
<pre><code>make build-run <i><*run-options...*></i></code></pre>

run already builded earlier:
<pre><code>make run-builded <i><*run-options...*></i></code></pre>

## Commands

Before each command you can see list of current options:  
  
![Alt text](/.documentation/images/options.png "List of options")  
  

1. **Show list of active tubes in queue**
```
list
```

2. **Change options**

- Changing single option:
```
set tube=$name
set delay=$value
set ttr=$value
set priority=$value
```

- Changing multiple options via a single set:  
<pre><code>set tube=$name delay=$value ttr=$value priority=$value  
set tube=$name priority=$value
<i>*etc...*</i></code></pre>

3. **Put data**

- From console
<pre><code>put -d <i>data...</i></code></pre>

- From file
<pre><code>put -f <i>path/to/file</i></code></pre>
