<script lang="ts">
	import * as d3 from 'd3';
	import { data } from '../lib/data';
	import { parseDataToTree } from '../lib/formatters';
	import type { D3Node } from '../lib/types';

	// Dimensions for the visualization
	let width = window.innerWidth;
	let height = window.innerHeight;

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
			.padding(1)
			.round(true)(root); // Apply the layout to the root
	}
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

<style>
	/* Ensure the treemap is responsive */
	svg {
		width: 100vw;
		height: 100vh;
	}

	text {
		pointer-events: none;
	}
</style>
