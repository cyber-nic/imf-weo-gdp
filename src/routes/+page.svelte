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

	$: d3data = parseDataToTree(data, selectedYear);

	// Create the treemap layout
	let treemapRoot: d3.HierarchyNode<D3Node>;

	$: {
		const root = d3
			.hierarchy(d3data)
			.sum((d) => d.value || 0) // Assign values only to leaf nodes
			.sort((a, b) => (b.value || 0) - (a.value || 0)); // Sort nodes by value

		treemapRoot = d3
			.treemap<D3Node>()
			.tile(d3.treemapBinary)
			.size([width, height])
			.padding(2)
			.round(true)(root); // Apply the layout to the root
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
</svg>

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

<style>
	/* Ensure the treemap is responsive */
	svg {
		width: 100vw;
		height: 100vh;
	}

	text {
		pointer-events: none;
	}

	#slider-container {
		position: fixed;
		bottom: 20px;
		left: 50%;
		transform: translateX(-50%);
		text-align: center;
		background: rgba(255, 255, 255, 0.8);
		padding: 10px;
		border-radius: 8px;
		box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
	}

	#year-slider {
		width: 300px;
		margin-top: 5px;
	}
</style>
