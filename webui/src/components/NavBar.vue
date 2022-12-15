<script>
export default {
	props: ['msg'],
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
	<div class="navbar d-flex justify-content-evenly align-items-center">
		<div>
            <i class="icon bi bi-house"></i>
        </div>
        <div>
            <i class="icon bi bi-search"></i>
        </div>
        <div>
            <i class="icon bi bi-person"></i>
        </div>
	</div>
</template>

<style>
    .icon {
        color: white;
        font-size: 2rem;
    }
    .navbar {
		background-color: #212121;
		overflow: hidden;
		position: fixed;
		bottom: 0;
		width: 75%;
		height: 75px;
	}
</style>
