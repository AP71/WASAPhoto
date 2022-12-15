<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			username: "",
			identifier: "",
		}
	},
	methods: {
		async doLogin() {
			this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.post("/session", {username: this.username,});
				this.identifier = response.data.identifier;
			} catch(e) {
				this.errormsg = e.toString();
			}
			this.loading = false
			if (this.identifier != "") {
				this.$router.push({ path: '/home', params: { identifier: this.identifier}});
			}
			
		}
	}
}
</script>

<template>
	<div class="d-flex min-vh-100 w-100 justify-content-center align-items-center" style="background-color: #383838">
		<div class="d-flex h-50 w-50 flex-column justify-content-center align-items-center rounded rounded-5" style="background-color: #212121;">
			<div class="fs-2 fw-bolder text-white pt-5">
				Login
			</div>
			<div class="pt-5">
				<input class="border border-3 border-success rounded-pill min-vh-25 min-vw-25 fs-4" v-model="username" placeholder="Username"/>
			</div>
			<div v-if="!loading" class="pt-5 pb-5">
				<button type="button" class="btn btn-outline-success rounded-pill fs-4" :disabled="username.length<3 || username.length>16" style="width: 150px" @click="doLogin">Login</button>
			</div>			
		</div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"/>
		<LoadingSpinner v-if="loading"/>
	</div>	
</template>

<style>
</style>
