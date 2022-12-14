<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			username: "",
		}
	},
	methods: {
		doLogin() {
            console.log(username)
        }
	},
}
</script>

<template>
	<div>
        <input type="text" v-model="username" placeholder="Username"/>
        <button class="button" @click="doLogin" :disabled="username.length<3 || username.length>16">Login</button>
    </div>
</template>

<style>
</style>
