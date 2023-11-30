const loader = new THREE.BufferGeometryLoader();

const Cards = (cfg) => {
    const map = {rToId:{}, idToR: {}}
    const _cards = new THREE.Group();
    const materials = new Material();
    const raycaster = new THREE.Raycaster();
    const mouse = new THREE.Vector2( 1, 1 );

    materials.MakeCustomMaterials(cfg.cards)

    function makeCardMesh(geometry, material) {
        const m = new THREE.Mesh(geometry, material)
        m.castShadow = true
        m.receiveShadow = true
        m.known = false
    
        return m
    }

    function makeGameCards() {
        loader.load(cfg.geometry, geometry => {

            geometry.computeVertexNormals();
            geometry.scale( 1, 0.1, 1 );
    
            const gridSizeX = cfg.grid[0] || 6
            const gridSizeZ = cfg.grid[1] || gridSizeX
            const interval = 2.1
    
            var n = 0
            for (let i = 0; i < gridSizeX; i++) {
                for (let j = 0; j < gridSizeZ; j++) {
                    const c = makeCardMesh(geometry, materials.Get('BaseMaterial'))
                    c.position.x = i * interval
                    c.position.z = j * interval
    
                    _cards.add(c)

                    map.idToR[c.id] = cfg.cardsIds[n]
                    map.rToId[cfg.cardsIds[n]] = c.id

                    n++
                }
            }
    
            cfg.scene.add(_cards)
    
            cfg.camera.Center(_cards)
        } );
    }
    makeGameCards();

    class Cards {
        id(card) {
            return map.idToR[card.id]
        }

        raycasterEvent(e) {
            e.preventDefault();

            mouse.x = ( e.clientX / window.innerWidth ) * 2 - 1;
            mouse.y = - ( e.clientY / window.innerHeight ) * 2 + 1;
            
            raycaster.setFromCamera( mouse, cfg.camera.Get() );
        
            const intersection = raycaster.intersectObjects( _cards.children );
        
            return {
                do: (fn) => { 
                    if (intersection.length > 0) {
                        fn(intersection[0].object)
                    } else {
                        fn(false)
                    }
                },
            }
        }

        Turn(reference_id, str = '', delay = 0, callback = () => {}) {
            const idx = map.rToId[reference_id]
            const card = _cards.getObjectById(idx)
            
            if (!card.known) {
                card.material = materials.Get(str)
                card.known = true
            }
        
            const tl = gsap.timeline({ onComplete: callback });
            
            const rZ = card.rotation.z == 0 ? Math.PI-0.00001 : 0
        
            tl.to(card.position, {y: 1.5, duration: 0.2, delay: delay})
                .to(card.rotation, {z: rZ, duration: 0.3})
                .to(card.position, {y: 0, duration: 0.1})
        }

        // remove paired
        RemovePaired(idA, idB, position, callback) {
            const idxA = map.rToId[idA]
            const cardA = _cards.getObjectById(idxA)
            const idxB = map.rToId[idB]
            const cardB = _cards.getObjectById(idxB)

            const tl = gsap.timeline({ onComplete: function() {
                cardA.visible = false;
                cardB.visible = false;
                
                callback()
            } });

            const delta = (Math.abs(position.x) + Math.abs(position.y)) / 2

            tl.to(cardA.position, {y: 1.5, duration: 0.1 * delta})
                .to(cardB.position, {y: 1.5, duration: 0.1 * delta}, '<')
                .to(cardA.position, {x: position.x, z: position.y, duration: 0.2 * delta})
                .to(cardB.position, {x: position.x, z: position.y, duration: 0.2 * delta}, '<')
                .to(cardA.scale, {x: 0.1, y: 0.01, z: 0.1, duration: 0.1 * delta})
                .to(cardB.scale, {x: 0.1, y: 0.01, z: 0.1, duration: 0.1 * delta}, '<')
        }
    }

    return new Cards
};