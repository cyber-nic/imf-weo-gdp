<script lang="ts">
	import { onMount } from 'svelte';
	import * as d3 from 'd3';
	import { data } from '../lib/data';
	import { parseDataToTree, getMinMaxYears } from '../lib/formatters';
	import type { D3Node } from '../lib/types';

	// Dimensions for the visualization
	let width = window.innerWidth;
	let height = window.innerHeight;

	const { minYear, maxYear } = getMinMaxYears(data[0]);

	// Define a color scale for diverse colors
	const colorScale = d3.scaleOrdinal(d3.schemeSet3);

	// Default to the current year
	const currentYear = new Date().getFullYear().toString();

	// Extract and filter data for the selected year
	let selectedYear = currentYear;

	// Visualization type: "treemap" or "pack"
	let visualizationType: 'treemap' | 'pack' = 'treemap';

	$: d3data = parseDataToTree(data, selectedYear);

	// Hierarchy root
	let root: d3.HierarchyNode<D3Node>;

	$: {
		root = d3
			.hierarchy(d3data)
			.sum((d) => d.value || 0)
			.sort((a, b) => (b.value || 0) - (a.value || 0));
	}

	// Treemap layout
	let treemapRoot: d3.HierarchyNode<D3Node>;
	$: if (visualizationType === 'treemap') {
		treemapRoot = d3
			.treemap<D3Node>()
			.tile(d3.treemapBinary)
			.size([width, height])
			.padding(2)
			.round(true)(root);
	}

	// Pack layout
	let packRoot: d3.HierarchyCircularNode<D3Node>;
	$: if (visualizationType === 'pack') {
		packRoot = d3.pack<D3Node>().size([width, height]).padding(3)(root);
	}

	// Handle window resize
	function updateDimensions() {
		width = window.innerWidth;
		height = window.innerHeight;
	}

	onMount(() => {
		window.addEventListener('resize', updateDimensions);
		return () => window.removeEventListener('resize', updateDimensions);
	});
</script>

<svg viewBox={`0 0 ${width} ${height}`} style="max-width: 100%; height: auto;">
	{#if visualizationType === 'treemap'}
		{#each treemapRoot.leaves() as leaf}
			<g transform={`translate(${leaf.x0},${leaf.y0})`}>
				<rect
					width={leaf.x1 - leaf.x0}
					height={leaf.y1 - leaf.y0}
					fill={colorScale(leaf.data.name)}
				/>
				<text
					x={(leaf.x1 - leaf.x0) / 2}
					y={(leaf.y1 - leaf.y0) / 2}
					dy="0.35em"
					text-anchor="middle"
					style={`pointer-events: none; font-size: 10px; fill: ${
						leaf.value && leaf.value / (treemapRoot.value || 1) > 0.5 ? 'white' : 'black'
					};`}
				>
					{leaf.data.name}
				</text>
			</g>
		{/each}
	{:else if visualizationType === 'pack'}
		{#each packRoot.descendants() as node}
			{#if node.depth > 0}
				<!-- Skip the root circle -->
				<g transform={`translate(${node.x},${node.y})`}>
					<circle r={node.r} fill={colorScale(node.data.name)} />
					<text
						dy="0.35em"
						text-anchor="middle"
						style={`pointer-events: none; font-size: ${node.r / 5}px; fill: black;`}
					>
						{node.data.name}
					</text>
				</g>
			{/if}
		{/each}
	{/if}
</svg>

<div id="controls-container">
	<div id="toggle-container">
		<select id="visualization-toggle" bind:value={visualizationType}>
			<option value="treemap">Treemap</option>
			<option value="pack">Pack</option>
		</select>
	</div>
	<div id="slider-container">
		<label for="year-slider">Year: {selectedYear}</label>
		<input
			id="year-slider"
			type="range"
			min={minYear}
			max={maxYear}
			value={selectedYear}
			on:input={(e) => (selectedYear = +e.target.value)}
		/>
	</div>
</div>

<style>
	/* Ensure the visualization is responsive */
	svg {
		width: 100vw;
		height: 100vh;
	}

	text {
		pointer-events: none;
	}

	#controls-container {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 20px;
		position: fixed;
		bottom: 20px;
		width: 100%;
	}

	#slider-container,
	#toggle-container {
		background: rgba(255, 255, 255, 0.8);
		padding: 10px;
		border-radius: 8px;
		box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
	}

	#year-slider {
		width: 300px;
		margin-top: 5px;
	}

	#visualization-toggle {
		margin-top: 5px;
	}
</style>
