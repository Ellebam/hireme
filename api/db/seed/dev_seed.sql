-- Development seed data
-- Run with: task db:seed

-- Create development user (matches AUTH_BYPASS_USER_ID)
INSERT INTO users (
    id,
    external_id,
    provider,
    email,
    email_verified,
    display_name,
    tier,
    cv_limit,
    storage_limit_bytes,
    storage_used_bytes,
    locale
) VALUES (
    'dev-user-001',
    'dev-bypass',
    'development',
    'dev@hireme.local',
    true,
    'Development User',
    'power',
    50,
    52428800,  -- 50MB
    0,
    'en'
) ON CONFLICT (id) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    tier = EXCLUDED.tier,
    cv_limit = EXCLUDED.cv_limit,
    storage_limit_bytes = EXCLUDED.storage_limit_bytes,
    updated_at = NOW();

-- Create sample CV for development
INSERT INTO cvs (
    id,
    user_id,
    title,
    schema_version,
    content,
    is_active
) VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'dev-user-001',
    'My Professional CV',
    '1.0.0',
    '{
        "schemaVersion": "1.0.0",
        "templateId": "modern",
        "locale": "en",
        "sections": [
            {
                "id": "sec-personal-001",
                "type": "personal",
                "order": 0,
                "visible": true,
                "content": {
                    "firstName": "Max",
                    "lastName": "Developer",
                    "jobTitle": "Senior Software Engineer",
                    "email": "max@example.com",
                    "phone": "+49 123 456789",
                    "location": "Berlin, Germany",
                    "links": [
                        {"type": "linkedin", "url": "https://linkedin.com/in/maxdev"},
                        {"type": "github", "url": "https://github.com/maxdev"}
                    ]
                }
            },
            {
                "id": "sec-summary-001",
                "type": "summary",
                "order": 1,
                "visible": true,
                "content": {
                    "text": "Experienced software engineer with 10+ years in building scalable web applications. Passionate about clean code, DevOps practices, and mentoring junior developers."
                }
            },
            {
                "id": "sec-experience-001",
                "type": "experience",
                "order": 2,
                "visible": true,
                "content": {
                    "entries": [
                        {
                            "id": "exp-001",
                            "company": "TechCorp GmbH",
                            "position": "Senior Software Engineer",
                            "location": "Berlin, Germany",
                            "startDate": "2020-01",
                            "endDate": null,
                            "current": true,
                            "description": "Leading development of cloud-native applications and microservices architecture.",
                            "highlights": [
                                "Architected microservices platform serving 1M+ users",
                                "Reduced deployment time by 70% through CI/CD improvements",
                                "Mentored team of 5 junior developers"
                            ]
                        },
                        {
                            "id": "exp-002",
                            "company": "StartupXYZ",
                            "position": "Full Stack Developer",
                            "location": "Munich, Germany",
                            "startDate": "2017-06",
                            "endDate": "2019-12",
                            "current": false,
                            "description": "Built and maintained multiple web applications from scratch.",
                            "highlights": [
                                "Developed MVP that secured Series A funding",
                                "Implemented real-time collaboration features"
                            ]
                        }
                    ]
                }
            },
            {
                "id": "sec-education-001",
                "type": "education",
                "order": 3,
                "visible": true,
                "content": {
                    "entries": [
                        {
                            "id": "edu-001",
                            "institution": "Technical University of Munich",
                            "degree": "M.Sc.",
                            "field": "Computer Science",
                            "location": "Munich, Germany",
                            "startDate": "2014-10",
                            "endDate": "2017-03",
                            "grade": "1.3"
                        }
                    ]
                }
            },
            {
                "id": "sec-skills-001",
                "type": "skills",
                "order": 4,
                "visible": true,
                "content": {
                    "categories": [
                        {
                            "id": "cat-001",
                            "name": "Programming Languages",
                            "skills": [
                                {"name": "Go", "level": "expert"},
                                {"name": "TypeScript", "level": "expert"},
                                {"name": "Python", "level": "advanced"}
                            ]
                        },
                        {
                            "id": "cat-002",
                            "name": "Frameworks & Tools",
                            "skills": [
                                {"name": "React", "level": "expert"},
                                {"name": "Next.js", "level": "advanced"},
                                {"name": "Docker", "level": "advanced"},
                                {"name": "Kubernetes", "level": "intermediate"}
                            ]
                        }
                    ]
                }
            },
            {
                "id": "sec-languages-001",
                "type": "languages",
                "order": 5,
                "visible": true,
                "content": {
                    "entries": [
                        {"id": "lang-001", "language": "German", "proficiency": "native"},
                        {"id": "lang-002", "language": "English", "proficiency": "fluent"}
                    ]
                }
            }
        ],
        "styling": {
            "primaryColor": "#2563eb",
            "fontFamily": "inter",
            "fontSize": "medium",
            "showIcons": true
        }
    }'::jsonb,
    true
) ON CONFLICT (id) DO UPDATE SET
    content = EXCLUDED.content,
    updated_at = NOW();

-- Create second sample CV (classic template, different content)
INSERT INTO cvs (
    id,
    user_id,
    title,
    schema_version,
    content,
    is_active
) VALUES (
    'b1ffbc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'dev-user-001',
    'Design Portfolio CV',
    '1.0.0',
    '{
        "schemaVersion": "1.0.0",
        "templateId": "classic",
        "locale": "en",
        "sections": [
            {
                "id": "sec-personal-002",
                "type": "personal",
                "order": 0,
                "visible": true,
                "content": {
                    "firstName": "Max",
                    "lastName": "Developer",
                    "jobTitle": "UX/UI Designer",
                    "email": "max.design@example.com",
                    "phone": "+49 123 456789",
                    "location": "Berlin, Germany",
                    "links": [
                        {"type": "portfolio", "url": "https://maxdesign.dev"},
                        {"type": "linkedin", "url": "https://linkedin.com/in/maxdesign"}
                    ]
                }
            },
            {
                "id": "sec-summary-002",
                "type": "summary",
                "order": 1,
                "visible": true,
                "content": {
                    "text": "Creative designer with a passion for user-centered design. 5+ years crafting intuitive interfaces for web and mobile applications."
                }
            },
            {
                "id": "sec-skills-002",
                "type": "skills",
                "order": 2,
                "visible": true,
                "content": {
                    "categories": [
                        {
                            "id": "cat-003",
                            "name": "Design Tools",
                            "skills": [
                                {"name": "Figma", "level": "expert"},
                                {"name": "Sketch", "level": "advanced"},
                                {"name": "Adobe XD", "level": "advanced"}
                            ]
                        },
                        {
                            "id": "cat-004",
                            "name": "Frontend",
                            "skills": [
                                {"name": "HTML/CSS", "level": "expert"},
                                {"name": "Tailwind CSS", "level": "advanced"},
                                {"name": "React", "level": "intermediate"}
                            ]
                        }
                    ]
                }
            }
        ],
        "styling": {
            "primaryColor": "#8b5cf6",
            "fontFamily": "inter",
            "fontSize": "medium",
            "showIcons": true
        }
    }'::jsonb,
    true
) ON CONFLICT (id) DO UPDATE SET
    content = EXCLUDED.content,
    updated_at = NOW();

-- Log completion
DO $$
BEGIN
    RAISE NOTICE '✅ Development seed data loaded successfully';
    RAISE NOTICE '   User: dev-user-001 (dev@hireme.local)';
    RAISE NOTICE '   CV 1: My Professional CV (modern)';
    RAISE NOTICE '   CV 2: Design Portfolio CV (classic)';
END $$;
