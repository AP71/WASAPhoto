<script setup>
import { RouterLink, RouterView } from 'vue-router'
import NavBar from './components/NavBar.vue'

</script>
<script>
export default {
	components: {
		NavBar 
	},
	data: function() {
		return {
			isLogged: false,
		}
	},
	methods: {

	},
	watch: {
    '$route': {
		handler () {
			this.isLogged = this.$profile.isLogged;
		},
    }
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
	<NavBar v-if="this.isLogged" />
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
