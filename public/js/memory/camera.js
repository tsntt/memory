const Camera = (() => {
    const cam = new THREE.PerspectiveCamera(40, window.innerWidth/window.innerHeight, 0.1, 1000);

    function pointCamera(mesh) {
        const boundingBox = new THREE.Box3();
        const center = new THREE.Vector3();

        boundingBox.setFromObject(mesh);
        boundingBox.getCenter(center);

        cam.position.copy(center);
        cam.position.y += boundingBox.getSize(new THREE.Vector3()).length();
        cam.lookAt(center);
    }

    class Camera {
        Center(mesh) {
            pointCamera(mesh)
        }

        Resize() {
            cam.aspect = window.innerWidth / window.innerHeight;
            cam.updateProjectionMatrix();
        }

        Get() {
            return cam
        }
    }

    return Camera
})();