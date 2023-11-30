const Lights = (scene) => {
    
    const light = new THREE.DirectionalLight( 0xffd17a, 1 );
    const target = new THREE.Object3D();
    target.position.set(1, 0, 8)
    light.position.set( 0, 20, 10);
    light.target = target;
    light.castShadow = true;
    light.shadow.camera.left = -15;
    light.shadow.camera.right = 15;
    light.shadow.camera.top = 15;
    light.shadow.camera.bottom = -15;
    light.shadow.mapSize.width = 4096; 
    light.shadow.mapSize.height = 4096;
    light.shadow.camera.near = 10;
    light.shadow.camera.far = 1000;
    scene.add(light);
    scene.add(target);


    const rectLight = new THREE.RectAreaLight( 0xff00ff, 10,  1.5, 1.5 );
    rectLight.position.set( 0, 0.1, 0 );
    rectLight.lookAt( 0, 0, 0 );
    rectLight.visible = false;
    scene.add( rectLight );

    const ambient = new THREE.AmbientLight( 0xffffff, 1.3 );
    scene.add( ambient );

    return [
        light,
        rectLight,
        ambient
    ]
}