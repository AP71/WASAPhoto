<script setup>
import NavBar from '../components/NavBar.vue'
</script>
<script>
export default {
	components: {
		NavBar
	},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			feed: null,
			identifier: "",
		}
	},
	methods: {
		
	}
}
</script>

<template>
	<div class="d-flex min-vh-100 w-100 justify-content-center align-items-center" style="background-color: #383838">
		<div class="min-vh-100 w-75" style="background-color: #2e2e2e;">		
			<div class="min-vh-75">
				Prova
			</div>
			<NavBar/>
		</div>
	</div>
</template>

<style>
</style>
