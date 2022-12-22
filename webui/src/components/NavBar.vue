<script>
export default {
	props: ['msg'],
	methods: {
		setStatus(i) {
			if (i==0) {
				this.$router.push({path: '/home'});
				document.getElementById("home").setAttribute('class', 'icon bi bi-house text-success');
				document.getElementById("search").setAttribute('class', 'icon bi bi-search text-white');
				document.getElementById("profile").setAttribute('class', 'icon bi bi-person text-white');
			}
			if (i==1) {
				this.$router.push({path: '/search'});
				document.getElementById("home").setAttribute('class', 'icon bi bi-house text-white');
				document.getElementById("search").setAttribute('class', 'icon bi bi-search text-success');
				document.getElementById("profile").setAttribute('class', 'icon bi bi-person text-white');
			}
			if (i==2) {
				this.$router.push({path: `/profiles/${this.$profile.username}`});
				document.getElementById("home").setAttribute('class', 'icon bi bi-house text-white');
				document.getElementById("search").setAttribute('class', 'icon bi bi-search text-white');
				document.getElementById("profile").setAttribute('class', 'icon bi bi-person text-success');
			}
		}
	}
}
</script>

<template>
	<div class="navbar">
		<div class="element" @click="setStatus(0)">
            <i id="home" class="icon bi bi-house"></i>
        </div>
        <div class="element" @click="setStatus(1)">
            <i id="search" class="icon bi bi-search"></i>
        </div>
        <div class="element" @click="setStatus(2)">
            <i id="profile" class="icon bi bi-person"></i>
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
		flex-grow: 1;
		justify-content: center;
		align-items: center;
	}
</style>
