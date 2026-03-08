<script>
  import { onMount } from 'svelte';

  let currentWorkspace = "alpha";
  let workspaces = ["alpha", "beta"];
  let agents = [
    { name: "ConversationAgent", status: "idle" },
    { name: "PlannerAgent", status: "active" }
  ];
  let messages = [];
  let prompt = "";

  function submitGoal() {
    if (!prompt) return;
    messages = [...messages, { role: "user", content: prompt }];
    prompt = "";
    // Call Wails Bridge: window.go.bridge.Bridge.SubmitGoal(currentWorkspace, prompt)
  }
</script>

<main class="h-screen flex flex-col bg-gray-900 text-white font-sans">
  <!-- Top Bar -->
  <header class="h-14 border-b border-gray-700 flex items-center px-4 justify-between">
    <div class="flex items-center space-x-4">
      <span class="font-bold tracking-widest text-blue-400">ORION</span>
      <select bind:value={currentWorkspace} class="bg-gray-800 border border-gray-600 rounded px-2 py-1 text-sm">
        {#each workspaces as ws}
          <option value={ws}>{ws}</option>
        {/each}
      </select>
    </div>
    <div class="flex space-x-2 text-xs">
      {#each agents as agent}
        <div class="flex items-center space-x-1 px-2 py-1 rounded bg-gray-800">
          <span class="w-2 h-2 rounded-full {agent.status === 'active' ? 'bg-green-500' : 'bg-gray-500'}"></span>
          <span>{agent.name}</span>
        </div>
      {/each}
    </div>
  </header>

  <!-- Main View -->
  <div class="flex-1 flex overflow-hidden">
    <!-- Chat Panel -->
    <section class="w-1/3 border-r border-gray-700 flex flex-col">
      <div class="flex-1 overflow-y-auto p-4 space-y-4">
        {#each messages as msg}
          <div class="p-3 rounded-lg {msg.role === 'user' ? 'bg-blue-900/30 border border-blue-700' : 'bg-gray-800 border border-gray-600'}">
            <p class="text-sm">{msg.content}</p>
          </div>
        {/each}
      </div>
      <div class="p-4 border-t border-gray-700">
        <div class="relative">
          <input
            bind:value={prompt}
            on:keydown={(e) => e.key === 'Enter' && submitGoal()}
            placeholder="Describe your goal..."
            class="w-full bg-gray-800 border border-gray-600 rounded-lg pl-4 pr-10 py-3 focus:outline-none focus:border-blue-500"
          />
          <button on:click={submitGoal} class="absolute right-3 top-3 text-blue-400">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
          </button>
        </div>
      </div>
    </section>

    <!-- Visualizer / Inspector -->
    <section class="flex-1 bg-black flex items-center justify-center">
      <div class="text-center text-gray-600">
        <p class="text-xl font-light italic">Cognitive Workspace Visualization</p>
        <p class="text-xs mt-2 uppercase tracking-widest">Select a node or goal to inspect</p>
      </div>
    </section>
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
  }
</style>
