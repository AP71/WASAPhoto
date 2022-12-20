<script>
export default {
	props: ['msg']
}
</script>

<template>
	<div v-if="!(this.msg==null)" class="alert alert-danger" role="alert">
		{{ this.msg }}
	</div>
</template>

<style>
</style>
