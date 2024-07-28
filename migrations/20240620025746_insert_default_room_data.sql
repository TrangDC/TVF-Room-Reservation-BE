INSERT INTO rooms (name, color, floor, office_id, description, image_url) VALUES
(
    'Tokyo',
    '#F72727',
    'Floor 1',
    (SELECT id FROM offices WHERE name='Hanoi'),
    'Capacity: 15-20 people\n\nFeatures:\n\nAir Conditioning: Central or high-capacity units\nProjector: High-resolution, with large screen\nAudio: Quality speakers and wireless microphones\nSeating: Comfortable chairs, arranged in U-shape or rectangular layout\nLighting: Adjustable LED ceiling lights, with curtains/blinds for natural light control\nAdditional: Whiteboard, high-speed Wi-Fi, multiple power outlets\nDesign: Professional and functional decor for a conducive meeting environment.',
    'https://ballantyneexecutivesuites.com/wp-content/uploads/2015/10/Depositphotos_13534536_original.jpg'
),
(
    'Singapore',
    '#27AFF7',
    'Floor 2',
    (SELECT id FROM offices WHERE name='Hanoi'),
    'Capacity: 10-15 people\n\nFeatures:\n\nAir Conditioning: Central or high-capacity units\nProjector: High-resolution, with large screen\nAudio: Quality speakers and wireless microphones\nSeating: Comfortable chairs, arranged in U-shape or rectangular layout\nLighting: Adjustable LED ceiling lights, with curtains/blinds for natural light control\nAdditional: Whiteboard, high-speed Wi-Fi, multiple power outlets\nDesign: Elegant decor for a productive meeting atmosphere.',
    'https://image-tc.galaxy.tf/wijpeg-ag5lo1h0y1bzpcp6qvh262qva/meeting-room-side-facing-warwick-san-francisco.jpg?width=2000'
),
(
    'London',
    '#37F727',
    'Floor 3',
    (SELECT id FROM offices WHERE name='Hanoi'),
    'Capacity: 20-25 people\n\nFeatures:\n\nAir Conditioning: Central or high-capacity units\nProjector: High-resolution, with large screen\nAudio: Quality speakers and wireless microphones\nSeating: Comfortable chairs, arranged in U-shape or rectangular layout\nLighting: Adjustable LED ceiling lights, with curtains/blinds for natural light control\nAdditional: Whiteboard, high-speed Wi-Fi, multiple power outlets\nDesign: Modern and sleek decor for a contemporary meeting environment.',
    'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTcOobDhZyyPizI535oy7D6JkffhAuJo5KU9g&s'
),
(
    'Sydney',
    '#F7EA27',
    'Floor 4',
    (SELECT id FROM offices WHERE name='Hanoi'),
    'Capacity: 20-25 people\n\nFeatures:\n\nAir Conditioning: Central or high-capacity units\nProjector: High-resolution, with large screen\nAudio: Quality speakers and wireless microphones\nSeating: Comfortable chairs, arranged in U-shape or rectangular layout\nLighting: Adjustable LED ceiling lights, with curtains/blinds for natural light control\nAdditional: Whiteboard, high-speed Wi-Fi, multiple power outlets\nDesign: Modern and sleek decor for a contemporary meeting environment.',
    'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTcOobDhZyyPizI535oy7D6JkffhAuJo5KU9g&s'
),
(
    'Da Nang - 1',
    '#BFBFBF',
    'Floor 5',
    (SELECT id FROM offices WHERE name='Da Nang'),
    'Capacity: 15-20 people\n\nFeatures:\n\nAir Conditioning: Central or high-capacity units\nProjector: High-resolution, with large screen\nAudio: Quality speakers and wireless microphones\nSeating: Comfortable chairs, arranged in U-shape or rectangular layout\nLighting: Adjustable LED ceiling lights, with curtains/blinds for natural light control\nAdditional: Whiteboard, high-speed Wi-Fi, multiple power outlets\nDesign: Professional and functional decor for a conducive meeting environment.',
    'https://image-tc.galaxy.tf/wijpeg-3q5inufelm3ghkidmo7qclyx2/1.jpg?width=2000'
),
(
    'Da Nang - 2',
    '#BFBFBF',
    'Floor 6',
    (SELECT id FROM offices WHERE name='Da Nang'),
    'Capacity: 15-20 people\n\nFeatures:\n\nAir Conditioning: Central or high-capacity units\nProjector: High-resolution, with large screen\nAudio: Quality speakers and wireless microphones\nSeating: Comfortable chairs, arranged in U-shape or rectangular layout\nLighting: Adjustable LED ceiling lights, with curtains/blinds for natural light control\nAdditional: Whiteboard, high-speed Wi-Fi, multiple power outlets\nDesign: Professional and functional decor for a conducive meeting environment.',
    'https://img.freepik.com/free-photo/room-used-official-event_23-2151054260.jpg?size=626&ext=jpg&ga=GA1.1.2116175301.1717804800&semt=ais_user'
),
(
    'HCM - 1',
    '#BFBFBF',
    'Floor 7',
    (SELECT id FROM offices WHERE name='Ho Chi Minh City'),
    'Capacity: 15-20 people\n\nFeatures:\n\nAir Conditioning: Central or high-capacity units\nProjector: High-resolution, with large screen\nAudio: Quality speakers and wireless microphones\nSeating: Comfortable chairs, arranged in U-shape or rectangular layout\nLighting: Adjustable LED ceiling lights, with curtains/blinds for natural light control\nAdditional: Whiteboard, high-speed Wi-Fi, multiple power outlets\nDesign: Professional and functional decor for a conducive meeting environment.',
    'https://tactic.au/static/7e3fc73a8a099c48a2d3f6f4c970c81b/meeting-room-design-2.jpg'
),
(
    'HCM - 2',
    '#BFBFBF',
    'Floor 8',
    (SELECT id FROM offices WHERE name='Ho Chi Minh City'),
    'Capacity: 15-20 people\n\nFeatures:\n\nAir Conditioning: Central or high-capacity units\nProjector: High-resolution, with large screen\nAudio: Quality speakers and wireless microphones\nSeating: Comfortable chairs, arranged in U-shape or rectangular layout\nLighting: Adjustable LED ceiling lights, with curtains/blinds for natural light control\nAdditional: Whiteboard, high-speed Wi-Fi, multiple power outlets\nDesign: Professional and functional decor for a conducive meeting environment.',
    'https://www.aver.com/Upload/Expert/31/Main.jpg'
),
(
    'HCM - 3',
    '#BFBFBF',
    'Floor 9',
    (SELECT id FROM offices WHERE name='Ho Chi Minh City'),
    'Capacity: 15-20 people\n\nFeatures:\n\nAir Conditioning: Central or high-capacity units\nProjector: High-resolution, with large screen\nAudio: Quality speakers and wireless microphones\nSeating: Comfortable chairs, arranged in U-shape or rectangular layout\nLighting: Adjustable LED ceiling lights, with curtains/blinds for natural light control\nAdditional: Whiteboard, high-speed Wi-Fi, multiple power outlets\nDesign: Professional and functional decor for a conducive meeting environment.',
    'https://www.groyalhotel.com.tw/wp-content/uploads/2020/09/%E6%9C%83%E8%AD%B0%E5%AE%A4201.jpg'
),
(
    'Japan - 1',
    '#BFBFBF',
    'Floor 10',
    (SELECT id FROM offices WHERE name='Japan'),
    'Capacity: 15-20 people\n\nFeatures:\n\nAir Conditioning: Central or high-capacity units\nProjector: High-resolution, with large screen\nAudio: Quality speakers and wireless microphones\nSeating: Comfortable chairs, arranged in U-shape or rectangular layout\nLighting: Adjustable LED ceiling lights, with curtains/blinds for natural light control\nAdditional: Whiteboard, high-speed Wi-Fi, multiple power outlets\nDesign: Professional and functional decor for a conducive meeting environment.',
    'https://cms.workplaceone.com/cr_images/meeting-rooms/_800x560_crop_center-center_none/sm-003_Mtg-Rm_WP1_26_Wellington_Street_E.jpg'
),
(
    'Japan - 2',
    '#BFBFBF',
    'Floor 11',
    (SELECT id FROM offices WHERE name='Japan'),
    'Capacity: 15-20 people\n\nFeatures:\n\nAir Conditioning: Central or high-capacity units\nProjector: High-resolution, with large screen\nAudio: Quality speakers and wireless microphones\nSeating: Comfortable chairs, arranged in U-shape or rectangular layout\nLighting: Adjustable LED ceiling lights, with curtains/blinds for natural light control\nAdditional: Whiteboard, high-speed Wi-Fi, multiple power outlets\nDesign: Professional and functional decor for a conducive meeting environment.',
    'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRKCIcIcSNDeuWv1_yA3rp1ajFVS0XTlh0u0johNm3ShgEGpJab7KzRL2CUr_8ckuCIKIE&usqp=CAU'
)