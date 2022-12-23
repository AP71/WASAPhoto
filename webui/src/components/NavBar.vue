<script>
export default {
	methods: {
		setStatus(i) {
			if (i==0) {
				this.$router.push({path: '/home'});
			}
			else if (i==1) {
				this.$router.push({path: '/search'});
			}
			else if (i==2) {
				this.$router.push({path: `/profiles/${this.$profile.username}`});
			}
		}
	}
}
</script>

<template>
	<div class="navbar">
		<div class="element" @click="setStatus(0)">
            <i class="icon bi bi-house"></i>
        </div>
        <div class="element" @click="setStatus(1)">
            <i class="icon bi bi-search"></i>
        </div>
        <div class="element" @click="setStatus(2)">
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
		display: flex;
		flex-direction: row;
		justify-content: space-evenly;
		align-items: center;
		background-color: #212121;
		overflow: hidden;
		position: fixed;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		width: 75%;
		height: 75px;
	}
	.element {
		display: flex;
		justify-content: center;
		align-items: center;
		min-width: 25%;
	}
</style>
