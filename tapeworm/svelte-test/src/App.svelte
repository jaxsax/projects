<script>
  import Tailwindcss from "./Tailwindcss.svelte";
  import { onMount } from 'svelte';
  import moment from "moment";
  import fuzzysort from 'fuzzysort';

  const strictsort = fuzzysort.new({threshold: -999});
  let links = [];
  let searchedLinks = [];
  let showRFC3339Time = false;
  let timer;
  let searchTerm = "";

  const debounce = v => {
		clearTimeout(timer);
		timer = setTimeout(() => {
			searchTerm = v;
		}, 100);
	}

  onMount(async() => {
    const res = await fetch(`/api/links`);
    links = await res.json();
    links = links.Links;
  })

  $: latestLink = () => {
    if (!links) {
      return null;
    }
    links[Math.floor(Math.random() * links.length)] || {};
  }

  $: {
    if (searchTerm) {
      strictsort.goAsync(searchTerm, links, { key: 'Title' })
        .then(p => {
          console.log(p);
          searchedLinks = p.map(x => x.obj);
        });
    } else {
      searchedLinks = links;
    }
  }

  function handleTimestampClick(_) {
    showRFC3339Time = !showRFC3339Time;
  }
</script>

<style>
.links {
  @apply mt-4;
}

.centerimage {
  margin: 0 auto;
}

.truncate {
  text-overflow: ellipsis;
}

/*
* {
  border: 1px solid #f00 !important;
} */
</style>

<main class="container mx-auto mt-8">
  <div>
    <img class="h-32 w-64 centerimage" src="/static/tapeworm-icon.png" alt="Tapeworm bot icon" />
    <input
      class="bg-white focus:outline-none focus:shadow-outline border
      border-gray-300 rounded-lg py-2 px-4 block w-full appearance-none
      leading-normal mt-2"
      on:keyup={({ target: { value }}) => debounce(value)}
      type="text"
      placeholder="{ latestLink ? latestLink.Title : "" }" />
    <span class="ml-2">
      { searchedLinks.length.toLocaleString() }
      result(s)
    </span>
  </div>
  <div class="links">
    {#each searchedLinks as link}
      <div class="my-4">
        <a class="font-semibold text-blue-600 visited:text-purple-600" href={link.Link}>{link.Title}</a>
        <span class=" border-b-2 border-dashed cursor-pointer" on:click={handleTimestampClick}>
          { showRFC3339Time ? moment.unix(link.CreatedTS).format() : moment.unix(link.CreatedTS).fromNow() }
        </span>
        <p class="text-xs truncate w-3/4">{link.Link}</p>
      </div>
    {/each}
  </div>
</main>
