<script setup>
import { RouterLink, RouterView } from 'vue-router'
import NavBar from './components/NavBar.vue'

</script>
<script>
export default {
	components: {
		NavBar 
	}
}
</script>

<template>
	<header class="header d-flex" style="height: 50px;background-color: #212121;">
		<div class="fs-1 fw-bolder text-white">
			WASA Photo
		</div>
	</header>
	<RouterView :key="this.$route.path"/>
	<NavBar v-if="this.$profile.isLogged" />
</template>

<style scoped>
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
