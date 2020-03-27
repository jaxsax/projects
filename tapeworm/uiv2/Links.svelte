<script>
    import moment from 'moment';

    let links = [];
    $: fetch(`http://192.168.1.114:9999/api/links`)
        .then(r => r.json())
        .then(data => {
            links = data.Links;
        })
</script>

<style>
    .link-title {
        font-size: 1.4em;
    }
    .link-item {
        margin: 1.2em 0;
    }
    .link-meta {
        margin-top: 0.4em;
    }
</style>

{#if links}
    <div class="ui relaxed list">
        {#each links as link}
            <div class="item link-item">
                <div class="content">
                    <div class="header">
                        <a href="{ link.Link }" rel="{ link.Title }" class='link-title'>
                            { link.Title }
                        </a>
                    </div>
                    <div class="description link-meta">
                        By { link.ExtraData.created_username } { moment.unix(link.CreatedTS).fromNow() }
                    </div>
                </div>
            </div>
        {:else}
            <p>No items :(</p>
        {/each}
    </div>
{:else}
    <p>loading...</p>
{/if}
