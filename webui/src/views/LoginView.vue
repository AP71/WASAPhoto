<script>
import LoadingSpinner from '../components/LoadingSpinner.vue'
import ErrorMsg from '../components/ErrorMsg.vue'

export default {
	components: {
		LoadingSpinner,
		ErrorMsg,
	},
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
				this.errormsg = e.response.data.message;
			}

			if (this.identifier != "") {
				this.$profile.setProfile(this.identifier, this.username);
				this.$axios.defaults.headers.common['Authorization'] = `Bearer ${this.identifier}`;
				this.$router.push({path: '/home'});
			}
			this.loading = false;
		},
	}
}
</script>

<template>
	<div class="d-flex min-vh-100 max-vh-100 w-100 justify-content-center align-items-center" style="background-color: #383838">
		<div class="d-flex h-50 w-50 flex-column justify-content-center align-items-center rounded rounded-5" style="background-color: #212121;">
			<div class="fs-2 fw-bolder text-white pt-5">
				Login
			</div>
			<div class="pt-5 pb-1">
				<input class="border border-3 border-success rounded-pill min-vh-25 min-vw-25 fs-4 text-indent" v-model="username" placeholder="Username"/>
			</div>
			<ErrorMsg v-if="this.errormsg" :msg="this.errormsg"/>
			<LoadingSpinner v-if="this.loading"/>
			<div v-if="!loading" class="pt-5 pb-5">
				<button type="button" class="btn btn-outline-success rounded-pill fs-4" :disabled="username.length<3 || username.length>16" style="width: 150px" @click="doLogin">Login</button>
			</div>			
		</div>
	</div>	
</template>

<style scoped>
	.text-indent {
	text-indent: 10px;
	}
</style>
