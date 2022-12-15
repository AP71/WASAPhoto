<script setup>
import { RouterLink, RouterView } from 'vue-router'
import LoginView from './views/LoginView.vue'
</script>
<script>
export default { components: { LoginView } }
</script>

<template>
	<header class="header d-flex" style="height: 50px;background-color: #212121;">
		<div class="fs-1 fw-bolder text-white">
			WASA Photo
		</div>
	</header>
	<RouterView />
</template>

<style>
	.header {
		justify-content: center;
		align-items: center;
		background-color: #212121;
		overflow: hidden;
		position: fixed;
		top: 0;
		width: 100%;
		height: 50px;
	}
</style>
