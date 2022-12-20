<script setup>
import { RouterLink, RouterView } from 'vue-router'
import NavBar from './components/NavBar.vue'

</script>
<script>
export default {
	components: {
		NavBar 
	},
	methods: {

	}
	
}
</script>

<template>
	<header class="header d-flex" style="height: 50px;background-color: #212121;">
		<div class="fs-1 fw-bolder text-white">
			WASA Photo
		</div>
	</header>
	<RouterView/>
	<div class="d-flex min-vh-100 w-100 justify-content-center align-items-center" style="background-color: #383838">
		<div class="min-vh-100 w-75" style="background-color: #2e2e2e;">		
			<div class="min-vh-75">
				<NavBar/>
			</div>
		</div>
	</div>
</template>

<style>
	.header {
		justify-content: center;
		align-items: center;
		background-color: #212121;
		overflow: hidden;
		position: fixed;
		z-index: 9999;
		top: 0;
		width: 100%;
		height: 50px;
	}
</style>
