<script lang="ts">
    export let className: string = "";
    export let title: string;
    export let description: string;
    export let canonical: string;
    export let schemas: string[] = [];
    export let facebook: { name: string, content: string}[] = [];
    export let twitter: { name: string, content: string}[] = [];

    $: scripts = schemas.map(
        (schema) => `
            <script type="application/ld+json">
                ${JSON.stringify(schema) + '<'}
            /script>
        `
    );
</script>

<svelte:head>
    <title>{title}</title>

    <meta name="description" content={description} />
    <link rel="canonical" href={canonical} />

    {#each facebook as { name, content }}
        <meta property={name} {content} />
    {/each}
    {#each twitter as { name, content }}
        <meta {name} {content} />
    {/each}

    {#each scripts as script}
        {@html script}
    {/each}
</svelte:head>


<div class="w-full mt-2 lg:mt-10 pb-24 {className}">
    <slot />
</div>