<script>
import LoginForm from '../components/LoginForm.vue';

export default {
    components: {
        LoginForm
	},
    data: function () {
        return {
            errormsg: null,
            loading: false,
        };
    },
    methods: {
        
    }
}
</script>

<template>
	<div>
		<LoginForm/>
	</div>
</template>

<style>
</style>
