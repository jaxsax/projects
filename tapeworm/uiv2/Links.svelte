<script>
    import moment from 'moment';

    let links = [];
    $: fetch(`http://192.168.1.114:9999/api/links`)
        .then(r => r.json())
        .then(data => {
            links = data.Links;
        })
</script>

{#if links}
    <div class="ui list">
        {#each links as link}
            <div class="item">
                <div class="content">
                    <div class="header">{ link.Title }</div>
                    <div class="description">
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
